package service

import (
	"../core"
)

type IndexService struct {
	*core.Service
}

func (c *IndexService) Index() *core.ViewResponse {
	c.Service.View.SetView("index.html")

	return &core.ViewResponse{View: c.Service.View}
}