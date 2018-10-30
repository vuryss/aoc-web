package core

type Response interface {
	GetBody() string
}

type ViewResponse struct {
	view *View
}


