package lib

import (
	"fmt"
	"release-notes-filler/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateReleaseNotes(task models.Task) {
	channel, err := CreateChannel(task.ID)

	if err != nil {
		message := fmt.Sprintf("Fatal: could not create broadcast channel: %v", err)
		logError(message, task, channel)
		return
	}

	logOperationStart("Starting updating release notes...", task, channel)

	logOperationStart("Finding the next app version...", task, channel)
	version, err := FetchEditableVersion(task.AppId)
	if err != nil {
		message := fmt.Sprintf("Failed to find the next app version (%v)", err)
		logError(message, task, channel)
		return
	}
	logOperationFinish(fmt.Sprintf("Found the next app version: %s", version.VersionString), task, channel)

	logOperationStart(fmt.Sprintf("Loading the list of enabled localizations of %s ...", version.VersionString), task, channel)
	localizations, err := FetchVersionLocalizations(version.Id)
	if err != nil {
		message := fmt.Sprintf("Failed to load the list of enabled localizations of %s (%v)", version.VersionString, err)
		logError(message, task, channel)
		return
	}
	logOperationFinish(fmt.Sprintf("Loaded the list of enabled localizations of %s", version.VersionString), task, channel)

	logOperationStart(fmt.Sprintf("Downloading translations from Lokalise (key = %s)", task.KeyName), task, channel)
	contents, err := FetchKeyContent(task.LokaliseProjectId, task.KeyName)
	if err != nil {
		message := fmt.Sprintf("Error downloading translations (%v)", err)
		logError(message, task, channel)
		return
	}
	logOperationFinish("Downloaded translations from Lokalise", task, channel)

	var updatedModels = map[string]string{}
	var failedModels []string
	for _, model := range localizations {
		var code = model.Locale
		if newCode, found := AppStoreLocaleToLokaliseLangaugeCode[model.Locale]; found {
			code = newCode
		}

		logOperationStart(fmt.Sprintf("Updating release notes for `%s`...", code), task, channel)

		var content = ""
		if newContent, found := contents[code]; found {
			content = newContent
		}

		if len(content) == 0 {
			logWarning(fmt.Sprintf("Unable to find contents for `%s`!", code), task, channel)
			logWarning(fmt.Sprintf("Skipped updating release notes for `%s`", code), task, channel)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModel, err := UpdateVersionLocalization(AppVersionLocalization{
			Id:           model.Id,
			Locale:       model.Locale,
			ReleaseNotes: content,
		})

		if err != nil {
			logWarning(fmt.Sprintf("Failed to update release notes for `%s`. Moving on.", code), task, channel)
			failedModels = append(failedModels, code)
			continue
		}

		updatedModels[code] = updatedModel.ReleaseNotes
		logOperationFinish(fmt.Sprintf("Updated release notes for `%s`.", code), task, channel)
	}

	var summaryBuilder strings.Builder
	summaryBuilder.WriteString("Summary\n\n")
	if len(updatedModels) > 0 {
		summaryBuilder.WriteString("✅ Updated\n\n")
		for code, text := range updatedModels {
			summaryBuilder.WriteString(fmt.Sprintf("%s:\n%s\n\n", code, text))
		}
	}

	if len(failedModels) > 0 {
		summaryBuilder.WriteString("❌ Not Updated\n")
		for _, code := range failedModels {
			summaryBuilder.WriteString(fmt.Sprintf("%s\n", code))
		}
	}

	logInfo(summaryBuilder.String(), task, channel)

	logOperationFinish("Completed updating release notes", task, channel)

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

func logWarning(message string, task models.Task, channel *Channel) {
	createLog(message, "warn", task, channel)
}

func logOperationStart(message string, task models.Task, channel *Channel) {
	createLog(message, "start", task, channel)
}

func logOperationFinish(message string, task models.Task, channel *Channel) {
	createLog(message, "finish", task, channel)
}

func logInfo(message string, task models.Task, channel *Channel) {
	createLog(message, "info", task, channel)
}

func createLog(message string, category string, task models.Task, channel *Channel) {
	var event = models.TaskEvent{
		TaskId:   task.ID,
		Category: category,
		Message:  message,
	}

	models.ModelStore.Create(&event)

	if channel != nil {
		data, err := event.AsJson()
		if err != nil {
			fmt.Fprintln(gin.DefaultErrorWriter, "goroutine - error serialize TaskEvent `%d` into JSON: %v", event.ID, err)
			return
		}
		channel.Broadcast <- data
	}
}
