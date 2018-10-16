package service

import (
	"net/http"
)

type Interface interface {
}

type Abstract struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	Parameters  map[string]string
}