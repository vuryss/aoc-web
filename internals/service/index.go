package service

import (
	"../core"
)

type IndexService struct {
	*core.Service
}

func (c *IndexService) Index() core.ViewResponse {
	c.Service.View.SetView("index.html")

	return c.Service.View.Response()
}