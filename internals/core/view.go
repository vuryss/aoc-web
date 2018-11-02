package core

import (
	"../helper"
	"bytes"
	"html/template"
	"log"
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
	viewInstance.viewData = make(map[string]interface{})

	return viewInstance
}

type View struct {
	config       *Config
	layoutFolder string
	layoutFile   string
	viewFolder	 string
	viewFile	 string
	viewData	 map[string]interface{}
}

func (view *View) SetView(file string) {
	view.viewFile = view.viewFolder + "/" + file
}

func (view *View) PassValue(name string, value interface{}) {
	view.viewData[name] = value
}

func (view *View) Render() string {
	templ := template.New("")
	view.addTemplateFunctions(templ)
	templ, err := templ.ParseFiles(view.layoutFile, view.viewFile)

	if err != nil {
		panic("VIEW: Render cannot read or parse views files, Error: " + err.Error())
	}

	b := bytes.Buffer{}
	log.Printf("View data: %v", view.viewData)
	err = templ.ExecuteTemplate(&b, path.Base(view.layoutFile), view.viewData)

	if err != nil {
		panic("VIEW: Render cannot render view, Error: " + err.Error())
	}

	return b.String()
}

func (view *View) addTemplateFunctions(t *template.Template) {
	t.Funcs(template.FuncMap{
		"N": func (n int) []struct{} { return make([]struct{}, n) },
		"int": func(a interface{}) int { return a.(int) },
	})
}
