package force_ssl

import (
	"crypto/tls"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strings"
)

type ForceSSLMiddleware struct {
	TrustXFPHeader     bool
	Enable301Redirects bool
	Message            string
}

func setDefaults(settings *ForceSSLMiddleware) {
	if settings.Message == "" {
		settings.Message = "SSL Required."
	}
}

func isNotSecure(secure *tls.ConnectionState, xfpHeader string, trustXfpHeader bool) bool {
	if trustXfpHeader {
		return xfpHeader != "https"
	} else {
		return secure == nil
	}
}

func (middleware *ForceSSLMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {
		setDefaults(middleware)

		secure := request.TLS
		xfpHeader := strings.ToLower(request.Header.Get("X-Forwarded-Proto"))

		if isNotSecure(secure, xfpHeader, middleware.TrustXFPHeader) {
			if middleware.Enable301Redirects {
				redirectURL := request.URL
				redirectURL.Scheme = "https"
				http.Redirect(
					writer.(http.ResponseWriter),
					request.Request,
					redirectURL.String(),
					http.StatusMovedPermanently,
				)
			} else {
				writer.WriteHeader(403)
				writer.(http.ResponseWriter).Write([]byte(middleware.Message))
			}
		} else {
			handler(writer, request)
		}
	}
}
