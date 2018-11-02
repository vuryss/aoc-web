package core

type Response interface {
	GetBody() string
}

type ViewResponse struct {
	View *View
}

func (vResponse *ViewResponse) GetBody() string {
	return vResponse.View.Render()
}
