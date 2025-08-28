package routes

import (
	"modbustohttp/pkg/views"
	"net/http"
)

var Routes = map[string]http.HandlerFunc{
	"/hello/": views.Hello,
}
