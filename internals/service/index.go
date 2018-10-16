package service

import "log"

type IndexService struct {
	*Abstract
}

func (service *IndexService) Index() {
	log.Print("This method is called")
}