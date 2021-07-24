package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type App struct {
	Name              string `yaml:"name"`
	Id                string `yaml:"app_id"`
	LokaliseProjectId string `yaml:"lokalise_project_id"`
	IconUrl           string `yaml:"icon_url"`
}

func FetchApps() []App {
	return _apps
}

func FetchApp(id string) (App, error) {
	if val, ok := _appsMap[id]; ok {
		return val, nil
	}

	message := fmt.Sprintf("App with id '%s' is not found", id)
	return App{}, errors.New(message)
}

// private methods

func readFromYaml() []App {
	content, err := ioutil.ReadFile("apps.yml")
	if err != nil {
		log.Fatalln(fmt.Printf("Fail to read apps.yml: %v", err))
	}

	var apps []App
	if err := yaml.Unmarshal(content, &apps); err != nil {
		log.Fatalln(fmt.Printf("Fail to unmarshall apps.yml: %v", err))
	}

	return apps
}

func buildMap(apps []App) map[string]App {
	var result = map[string]App{}

	for _, app := range apps {
		result[app.Id] = app
	}
	return result
}

var _apps []App = readFromYaml()
var _appsMap map[string]App = buildMap(_apps)
