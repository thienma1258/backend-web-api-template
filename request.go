package api

import (
	"context"
	"log"
	"net/http"
)

// HTTP headers
const (
	HeaderContentType = "Content-Type"
)

// Content type
const (
	JsonContentType = "Content-Type"
)

// InternalCode
const (
	InternalCodeSuccess = "0"
	InternalError       = "999"
)

// Success Message
const (
	SuccessMessage = "success"
)

type RequestContext struct {
	RemoteIP       string
	RemoteLocation string
	RemoteTimezone string
	Auth           *AuthUser
	ctx            context.Context
	idx            int
	middlewares    []func(*RequestContext)
	Writer         http.ResponseWriter
	Request        *http.Request
	Route          Router
}

type AuthUser struct {
	Username    string
	Token       string
	Roles       []string
	Permissions []string
}

type mapObject = map[string]interface{}

func (ctx *RequestContext) GetContext() context.Context {
	return ctx.ctx
}

func (ctx *RequestContext) Next() {
	if ctx.idx >= len(ctx.middlewares) {
		panic("end of the function chaining")
	}

	ctx.idx++
	ctx.middlewares[ctx.idx](ctx)
}

// SendJSON encodes data as JSON object and returns it to client
func (ctx *RequestContext) SendJSON(statusCode int, internalCode string, message string, data interface{}) error {
	return SendJSON(ctx.Writer, statusCode, internalCode, message, data)
}

// SendSuccess sends response with http status 200 and verdict success
func (ctx *RequestContext) SendSuccess(data interface{}) error {
	return ctx.SendJSON(http.StatusOK, InternalCodeSuccess, SuccessMessage, data)
}

// SendError sends internal error response to client
func (ctx *RequestContext) SendError(err error) error {
	return ctx.SendJSON(http.StatusInternalServerError, InternalError, err.Error(), mapObject{})
}

// SendJSON encodes data as JSON object and returns it to client
func SendJSON(
	w http.ResponseWriter,
	statusCode int,
	internalCode string,
	message string,
	data interface{},
) error {
	w.Header().Set(HeaderContentType, JsonContentType)
	w.WriteHeader(statusCode)

	obj := map[string]interface{}{
		"internalCode": internalCode,
		"message":      message,
		"data":         data,
		"time":         GetCurrentISOTime(),
	}

	body, err := json.Marshal(obj)
	if err != nil {
		log.Printf("cannot marshal response data: %v", err)
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("cannot write response data: %v", err)
		return err
	}

	return nil
}
