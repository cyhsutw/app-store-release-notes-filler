package lib

import (
	"fmt"
	"release-notes-filler/models"
	"time"
)

func UpdateReleaseNotes(task models.Task) {
	logInfo("⚡️ Starting updating release notes...", task)

	logInfo("🌐 Finding the next app version...", task)
	version, err := FetchEditableVersion(task.AppId)
	if err != nil {
		message := fmt.Sprintf("❌ Failed to find the next app version (%v)", err)
		logError(message, task)
		return
	}
	logInfo(fmt.Sprintf("✅ Found the next app version: %s", version.VersionString), task)

	logInfo(fmt.Sprintf("🌐 Loading the list of enabled localizations of %s ...", version.VersionString), task)
	localizations, err := FetchVersionLocalizations(version.Id)
	if err != nil {
		message := fmt.Sprintf("❌ Failed to load the list of enabled localizations of %s (%v)", version.VersionString, err)
		logError(message, task)
		return
	}
	logInfo(fmt.Sprintf("✅ Loaded the list of enabled localizations of %s", version.VersionString), task)

	logInfo(fmt.Sprintf("💾 Downloading translations from Lokalise (key = %s)", task.KeyName), task)
	contents, err := FetchKeyContent(task.LokaliseProjectId, task.KeyName)
	if err != nil {
		message := fmt.Sprintf("❌ Error downloading translations (%v)", err)
		logError(message, task)
		return
	}
	logInfo("✅ Downloaded translations from Lokalise", task)

	var updatedModels = map[string]string{}
	var failedModels []string
	for _, model := range localizations {
		var code = model.Locale
		if newCode, found := AppStoreLocaleToLokaliseLangaugeCode[model.Locale]; found {
			code = newCode
		}

		logInfo(fmt.Sprintf("📝 Updating release notes for %s ...", code), task)

		var content = ""
		if newContent, found := contents[code]; found {
			content = newContent
		}

		if len(content) == 0 {
			logInfo(fmt.Sprintf("⚠️ Unable to find contents for %s!", code), task)
			logInfo(fmt.Sprintf("✂️ Skipped updating release notes for %s", code), task)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModel, err := UpdateVersionLocalization(AppVersionLocalization{
			Id:           model.Id,
			Locale:       model.Locale,
			ReleaseNotes: content,
		})

		if err != nil {
			logInfo(fmt.Sprintf("❌ Failed to update release notes for %s. Moving on.", code), task)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModels[code] = updatedModel.ReleaseNotes
		logInfo(fmt.Sprintf("✅ Updated release notes for %s.", code), task)
	}

	logInfo("------ Summary ------", task)

	if len(updatedModels) > 0 {
		logInfo("✅ Updated", task)
		for code, text := range updatedModels {
			logInfo(fmt.Sprintf("%10s → %s", code, text), task)
		}
	}

	if len(failedModels) > 0 {
		logInfo("❌ Not Updated", task)
		for _, code := range failedModels {
			logInfo(fmt.Sprintf("%10s", code), task)
		}
	}

	logInfo("👌 Completed updating release notes", task)

	models.ModelStore.Model(&task).Updates(models.Task{
		Status:      "succeeded",
		CompletedAt: time.Now(),
	})
}

func logError(message string, task models.Task) {
	createLog(message, "error", task)
	models.ModelStore.Model(&task).Updates(models.Task{
		Status:      "failed",
		CompletedAt: time.Now(),
	})
}

func logInfo(message string, task models.Task) {
	createLog(message, "info", task)
}

func createLog(message string, logType string, task models.Task) {
	models.ModelStore.Create(&models.TaskLog{
		TaskId:  task.ID,
		LogType: logType,
		Message: message,
	})
}
