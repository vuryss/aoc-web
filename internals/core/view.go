package core

import (
	"../helper"
	"io/ioutil"
	"strings"
)

func NewView(config *Config) *View {
	viewInstance := &View{config: config}

	// Set views directory
	viewDirectory, found := config.GetString("views_directory")

	if !found {
		panic("VIEW: Cannot find views directory in configuration!")
	}

	viewInstance.viewFolder = helper.ResolveProjectFile(viewDirectory)

	// Set layouts directory
	layoutsDirectory, found := config.GetString("layouts_directory")

	if !found {
		panic("VIEW: Cannot find layouts directory in configuration!")
	}

	viewInstance.layoutFolder = helper.ResolveProjectFile(layoutsDirectory)

	// Set default layout file
	defaultLayout, found := config.GetString("default_layout")

	if !found {
		defaultLayout = "default.html"
	}

	viewInstance.layoutFile = helper.ResolveProjectFile(layoutsDirectory + "/" + defaultLayout)

	return viewInstance
}

type View struct {
	config       *Config
	layoutFolder string
	layoutFile   string
	loadLayout	 bool
	viewFolder	 string
	viewFile	 string
	isDisabled   bool
}

func (view *View) SetView(file string) {
	view.viewFile = view.viewFolder + "/" + file
}

func (view *View) Render() string {
	// Read layout file
	layoutContent, err := ioutil.ReadFile(view.layoutFile)

	if err != nil {
		panic("VIEW: Render cannot read layout file: " + view.layoutFile)
	}

	// Read view file
	viewContent, err := ioutil.ReadFile(view.viewFile)

	if err != nil {
		panic("VIEW: Render cannot read view file: " + view.viewFile)
	}

	// Replace special placeholder in layout with view content
	response := strings.Replace(string(layoutContent), "{{__VIEW__}}", string(viewContent), 1)

	return response
}
