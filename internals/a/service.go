package a

import (
	"net/http"
)

type Service struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	Parameters  map[string]string
	View        View
}