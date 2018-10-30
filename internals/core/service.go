package core

import (
	"net/http"
)

type Service struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	Parameters  map[string]string
	Config 		*Config
	View        *View
}