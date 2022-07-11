package api

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

//Middleware middleware for the http request
type Middleware interface {
	Handle(ctx *RequestContext)
}

//MonitoringMiddleware parses parameters middleware
type MonitoringMiddleware struct {
	LogMode LOG_MODE
}

// Handle parses parameters and saving to the request context based on request type
func (m *MonitoringMiddleware) Handle(ctx *RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s: %s", err, debug.Stack())
			log.Printf("Recovered in service %+v", err)
			return
		}
	}()
	if m.LogMode.isDebugMode() || m.LogMode.isWarningMode() {
		startTime := time.Now().UnixNano()
		r := ctx.Request
		w := ctx.Writer
		code := r.Header.Get("CF-Ipcountry")
		if len(code) > 0 {
			log.Printf("%s %s \t%dms \t%sb \t: %s\t--> %s", code, GetIPAddresses(ctx.Request), (time.Now().UnixNano()-startTime)/1000000, ctx.Writer.Header().Get("Expected-Size"), r.Method, "https://"+r.Host+r.URL.Path+"?"+r.URL.RawQuery)
		} else {
			log.Printf("%s \t%dms \t%sb \t: %s\t--> %s", GetIPAddresses(ctx.Request), (time.Now().UnixNano()-startTime)/1000000, w.Header().Get("Expected-Size"), r.Method, "https://"+r.Host+r.URL.Path+"?"+r.URL.RawQuery)
		}
	}

	ctx.Next()
}

//ParseInfoMiddleware parses parameters middleware
type ParseInfoMiddleware struct{}

// Handle parses parameters and saving to the request context based on request type
func (m *ParseInfoMiddleware) Handle(ctx *RequestContext) {

	ipAddress := GetIPAddresses(ctx.Request)
	ctx.RemoteIP = ipAddress
	ctx.Next()
}

// CorsMiddleware filter middleware
type CorsMiddleware struct {
	CORS *CORS
}

// Handle filter by method and add some header flags to response
func (m *CorsMiddleware) Handle(ctx *RequestContext) {
	r := ctx.Request
	w := ctx.Writer

	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s: %s", err, debug.Stack())
			log.Printf("Recovered in service %+v", err)
			//utils.ResponseError(utils.ERROR_UNKNOWN_ERROR, w)
		}
	}()

	// set response headers
	w.Header().Set("Access-Control-Allow-Origin", m.CORS.AllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", m.CORS.AllowMethods)
	w.Header().Set("Access-Control-Allow-Headers", m.CORS.AllowHeaders)
	// preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != ctx.Route.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx.Next()
}
