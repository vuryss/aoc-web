package core

import "../helper"

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

	viewInstance.layoutFile = defaultLayout

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
	view.viewFile = helper.ResolveProjectFile(view.viewFolder + "/" + file)
}

func (view *View) Render() string {
	return ""
}

func (view *View) Response() ViewResponse {
	return ViewResponse{view: view}
}
