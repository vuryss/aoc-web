package service

import (
	"../../internals/a"
	"log"
)

type IndexService struct {
	*a.Service
}

func (service *IndexService) Index() {
	log.Print("This method is called")
}