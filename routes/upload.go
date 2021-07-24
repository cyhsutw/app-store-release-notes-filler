package routes

import (
	"fmt"
	"log"
	"net/http"
	"release-notes-filler/lib"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	appId := c.DefaultPostForm("app_id", "")
	if len(appId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "app not provided"})
		return
	}

	keyName := c.DefaultPostForm("key_name", "")
	if len(keyName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key not provided"})
		return
	}

	app, err := lib.FetchApp(appId)
	if err != nil {
		message := fmt.Sprintf("fetch app error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	contents, err := lib.FetchKeyContent(app.LokaliseProjectId, keyName)

	if err != nil {
		message := fmt.Sprintf("fetch lokalise key '%s' error: %v", keyName, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": message})
		return
	}

	version, err := lib.FetchEditableVersion(app.Id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusAccepted)
		return
	}

	localizations, err := lib.FetchVersionLocalizations(version.Id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusAccepted)
		return
	}

	for _, model := range localizations {
		var code = model.Locale
		if newCode, found := lib.AppStoreLocaleToLokaliseLangaugeCode[model.Locale]; found {
			code = newCode
		}

		var content = ""
		if newContent, found := contents[code]; found {
			content = newContent
		}

		if len(content) == 0 {
			log.Fatalln("data not found")
		}

		model.ReleaseNotes = content
		updatedModel, err := lib.UpdateVersionLocalization(model)

		if err != nil {
			log.Printf("err update %s: %v", model.Locale, err)
		}

		log.Printf("updated -> %s", updatedModel.ReleaseNotes)
	}

	c.Status(http.StatusCreated)
}
