package force_ssl

import (
	"net/http"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

type JSON map[string]interface{}

func simpleEndpoint(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(JSON{
		"email":    "edward@example.com",
		"password": "password",
	})
}

func NewAPI() http.Handler {
	api := rest.NewApi()
	api.Use(&ForceSSLMiddleware{})
	simpleEndpoint := rest.AppSimple(simpleEndpoint)
	api.SetApp(simpleEndpoint)
	return api.MakeHandler()
}

func TestUnconfiguredForceSSLMiddleware(t *testing.T) {
	handler := NewAPI()
	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recorded := test.RunRequest(t, handler, req)
	recorded.CodeIs(http.StatusForbidden)
	recorded.BodyIs("SSL Required.")
}
