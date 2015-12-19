package forceSSL

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strings"
)

// Middleware is a go-json-rest middleware. It requires all
// requests to a go-json-rest server to be over SSL.
type Middleware struct {
	// Trust X-Forwarded-Proto headers (this could allow a client
	// to spoof whether they were using HTTPS).
	// Optional, defaults to false.
	TrustXFPHeader bool

	// Enables 301 redirects to the HTTPS version of the request.
	// Optional, defaults to false.
	Enable301Redirects bool

	// Allows a custom response message when forcing SSL without redirect.
	// Optional, defaults to "SSL Required."
	Message string
}

func setDefaults(settings *Middleware) {
	if settings.Message == "" {
		settings.Message = "SSL Required."
	}
}

func isNotSecure(request *rest.Request, xfpHeader string, trustXfpHeader bool) bool {
	if trustXfpHeader {
		return xfpHeader != "https"
	}

	return request.TLS == nil && request.URL.Scheme != "https"
}

// MiddlewareFunc makes forceSSL.Middleware implement the rest.Middleware interface.
func (middleware *Middleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {
		setDefaults(middleware)
		goHTTPWriter := writer.(http.ResponseWriter)

		xfpHeader := strings.ToLower(request.Header.Get("X-Forwarded-Proto"))

		if isNotSecure(request, xfpHeader, middleware.TrustXFPHeader) {
			if middleware.Enable301Redirects {
				redirectURL := request.URL
				redirectURL.Scheme = "https"
				http.Redirect(
					goHTTPWriter,
					request.Request,
					redirectURL.String(),
					http.StatusMovedPermanently,
				)
			} else {
				writer.WriteHeader(http.StatusForbidden)
				goHTTPWriter.Write([]byte(middleware.Message))
			}
		} else {
			handler(writer, request)
		}
	}
}
