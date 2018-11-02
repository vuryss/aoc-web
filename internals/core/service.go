package core

import (
	"net/http"
	"reflect"
)

type Service struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	Parameters  map[string]string
	Config 		*Config
	View        *View
}

func (s *Service) RedirectWithViewData(service string, action string) Response {
	a, exists := ServicesList[service]

	if !exists {
		panic("Redirect to invalid service")
	}

	reflectInstance := reflect.New(reflect.TypeOf(a).Elem())

	abstractService := reflect.Indirect(reflectInstance).FieldByName("Service")
	abstractService.Set(reflect.ValueOf(&Service{
		Request		: s.Request,
		Response	: s.Response,
		Parameters	: map[string]string{},
		Config		: s.Config,
		View        : s.View,
	}))

	methodRef := reflectInstance.MethodByName(action)
	responseValue := methodRef.Call(nil)[0]

	return responseValue.Interface().(Response)
}
