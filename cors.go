package api

import (
	"fmt"
	"strings"
)

// CORS contains headers to control Cross-Origin Resource Sharing (CORS)
type CORS struct {
	AllowOrigin      string
	AllowCredentials string
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	XFrameOptions    string
}

// AllOrigins allows request from all origins
const AllOrigins = "*"

// DefaultCORS contains default CORS values
var DefaultCORS = &CORS{
	AllowOrigin:      AllOrigins,
	AllowCredentials: "true",
	AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE",
	AllowHeaders:     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token",
	ExposeHeaders:    "*",
}

// SetAllowOriginHeader sets the allow origin header value
func (cors *CORS) SetAllowOriginHeader(origin string) error {
	origin = strings.TrimSpace(origin)

	if origin != "" && origin != AllOrigins && !strings.HasPrefix(origin, "http://") &&
		!strings.HasPrefix(origin, "https://") {
		return fmt.Errorf("invalid origin %s", origin)
	}

	cors.AllowOrigin = origin
	return nil
}
