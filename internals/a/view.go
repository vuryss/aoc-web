package a

import "../helper"

type View struct {
	initialized  bool
	layoutFolder string
	layoutFile   string
	loadLayout	 bool
	viewFolder	 string
	viewFile	 string
}

func (view *View) SetView(path string) {
	view.ensureView()
}

func (view *View) ensureView() {
	if !view.initialized {
		view.viewFolder = helper.Config.Get("view_folder")
	}
}