package lib

import (
	"fmt"
	"release-notes-filler/models"
	"time"
)

func UpdateReleaseNotes(task models.Task) {
	channel, err := CreateChannel(task.ID)

	if err != nil {
		message := fmt.Sprintf("Fatal: could not create broadcast channel: %v", err)
		logError(message, task, channel)
		return
	}

	logInfo("⚡️ Starting updating release notes...", task, channel)

	logInfo("🌐 Finding the next app version...", task, channel)
	version, err := FetchEditableVersion(task.AppId)
	if err != nil {
		message := fmt.Sprintf("❌ Failed to find the next app version (%v)", err)
		logError(message, task, channel)
		return
	}
	logInfo(fmt.Sprintf("✅ Found the next app version: %s", version.VersionString), task, channel)

	logInfo(fmt.Sprintf("🌐 Loading the list of enabled localizations of %s ...", version.VersionString), task, channel)
	localizations, err := FetchVersionLocalizations(version.Id)
	if err != nil {
		message := fmt.Sprintf("❌ Failed to load the list of enabled localizations of %s (%v)", version.VersionString, err)
		logError(message, task, channel)
		return
	}
	logInfo(fmt.Sprintf("✅ Loaded the list of enabled localizations of %s", version.VersionString), task, channel)

	logInfo(fmt.Sprintf("💾 Downloading translations from Lokalise (key = %s)", task.KeyName), task, channel)
	contents, err := FetchKeyContent(task.LokaliseProjectId, task.KeyName)
	if err != nil {
		message := fmt.Sprintf("❌ Error downloading translations (%v)", err)
		logError(message, task, channel)
		return
	}
	logInfo("✅ Downloaded translations from Lokalise", task, channel)

	var updatedModels = map[string]string{}
	var failedModels []string
	for _, model := range localizations {
		var code = model.Locale
		if newCode, found := AppStoreLocaleToLokaliseLangaugeCode[model.Locale]; found {
			code = newCode
		}

		logInfo(fmt.Sprintf("📝 Updating release notes for %s ...", code), task, channel)

		var content = ""
		if newContent, found := contents[code]; found {
			content = newContent
		}

		if len(content) == 0 {
			logInfo(fmt.Sprintf("⚠️ Unable to find contents for %s!", code), task, channel)
			logInfo(fmt.Sprintf("✂️ Skipped updating release notes for %s", code), task, channel)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModel, err := UpdateVersionLocalization(AppVersionLocalization{
			Id:           model.Id,
			Locale:       model.Locale,
			ReleaseNotes: content,
		})

		if err != nil {
			logInfo(fmt.Sprintf("❌ Failed to update release notes for %s. Moving on.", code), task, channel)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModels[code] = updatedModel.ReleaseNotes
		logInfo(fmt.Sprintf("✅ Updated release notes for %s.", code), task, channel)
	}

	logInfo("------ Summary ------", task, channel)

	if len(updatedModels) > 0 {
		logInfo("✅ Updated", task, channel)
		for code, text := range updatedModels {
			logInfo(fmt.Sprintf("%s:\n%s", code, text), task, channel)
		}
	}

	if len(failedModels) > 0 {
		logInfo("❌ Not Updated", task, channel)
		for _, code := range failedModels {
			logInfo(fmt.Sprintf("%10s", code), task, channel)
		}
	}

	logInfo("👌 Completed updating release notes", task, channel)

	timestamp := time.Now()
	models.ModelStore.Model(&task).Updates(models.Task{
		Status:      "succeeded",
		CompletedAt: &timestamp,
	})

	DestroyChannel(task.ID)
}

func logError(message string, task models.Task, channel *Channel) {
	createLog(message, "error", task, channel)

	timestamp := time.Now()
	models.ModelStore.Model(&task).Updates(models.Task{
		Status:      "failed",
		CompletedAt: &timestamp,
	})
}

func logInfo(message string, task models.Task, channel *Channel) {
	createLog(message, "info", task, channel)
}

func createLog(message string, logType string, task models.Task, channel *Channel) {
	models.ModelStore.Create(&models.TaskLog{
		TaskId:  task.ID,
		LogType: logType,
		Message: message,
	})

	if channel != nil {
		channel.Broadcast <- message
	}
}
