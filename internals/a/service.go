package a

import (
	"../core"
	"net/http"
)

type Service struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	Parameters  map[string]string
	Config 		*core.Config
	View        View
}