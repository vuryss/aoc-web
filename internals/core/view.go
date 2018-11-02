package core

import (
	"../helper"
	"bytes"
	"html/template"
	"path"
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

	// Init view data
	viewInstance.viewData = make(map[string]string)

	return viewInstance
}

type View struct {
	config       *Config
	layoutFolder string
	layoutFile   string
	viewFolder	 string
	viewFile	 string
	viewData	 map[string]string
}

func (view *View) SetView(file string) {
	view.viewFile = view.viewFolder + "/" + file
}

func (view *View) PassValue(name, value string) {
	view.viewData[name] = value
}

func (view *View) Render() string {
	templ, err := template.ParseFiles(view.layoutFile, view.viewFile)

	if err != nil {
		panic("VIEW: Render cannot read or parse views files, Error: " + err.Error())
	}

	b := bytes.Buffer{}

	err = templ.ExecuteTemplate(&b, path.Base(view.layoutFile), view.viewData)

	if err != nil {
		panic("VIEW: Render cannot render view, Error: " + err.Error())
	}

	return b.String()
}
