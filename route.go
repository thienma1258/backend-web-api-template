package api

import "net/http"

var indexRoute = Router{
	Path:   "/",
	Method: http.MethodGet,
}

var pingRoute = Router{
	Path:   "/ping",
	Method: http.MethodGet,
}

type HttpMethod = string

type Router struct {
	Path   string
	Method HttpMethod
}

type RouterHandler struct {
	Route   Router
	Handler func(ctx *RequestContext)
}

type OpHandler = func(http.ResponseWriter, *http.Request)
