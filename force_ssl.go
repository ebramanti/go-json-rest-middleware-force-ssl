package forceSSL

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"net/url"
	"strings"
)

type ForceSSLMiddleware struct {
	TrustXFPHeader     bool
type Middleware struct {
	Enable301Redirects bool
	Message            string
}

func setDefaults(settings *Middleware) {
	if settings.Message == "" {
		settings.Message = "SSL Required."
	}
}

func isNotSecure(url *url.URL, xfpHeader string, trustXfpHeader bool) bool {
	if trustXfpHeader {
		return xfpHeader != "https"
	} else {
		return url.Scheme != "https"
	}
}

func (middleware *Middleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {
		setDefaults(middleware)
		goHTTPWriter := writer.(http.ResponseWriter)

		url := request.URL
		xfpHeader := strings.ToLower(request.Header.Get("X-Forwarded-Proto"))

		if isNotSecure(url, xfpHeader, middleware.TrustXFPHeader) {
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
				writer.WriteHeader(403)
				goHTTPWriter.Write([]byte(middleware.Message))
			}
		} else {
			handler(writer, request)
		}
	}
}
