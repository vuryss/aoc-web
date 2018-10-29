package service

import (
	"../../internals/a"
)

type IndexService struct {
	*a.Service
}

func (c *IndexService) Index() {
	c.Service.View.SetView("index.html")
}