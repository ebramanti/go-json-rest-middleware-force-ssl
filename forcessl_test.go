package forceSSL

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"net/http"
	"testing"
)

type JSON map[string]interface{}

var (
	simplePostData = JSON{
		"email":    "edward@example.com",
		"password": "password",
	}
)

func simpleGetEndpoint(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(simplePostData)
}

func simplePostEndpoint(w rest.ResponseWriter, r *rest.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	r.DecodeJsonPayload(&body)
	w.WriteJson(body)
}

func NewAPI(forceSSLMiddleware *Middleware) http.Handler {
	api := rest.NewApi()
	api.Use(forceSSLMiddleware)
	router, _ := rest.MakeRouter(
		rest.Post("/", simplePostEndpoint),
		rest.Get("/", simpleGetEndpoint),
	)
	api.SetApp(router)
	return api.MakeHandler()
}

func TestUnconfiguredForceSSLMiddleware(t *testing.T) {
	handler := NewAPI(&Middleware{})

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recorded := test.RunRequest(t, handler, req)

	recorded.CodeIs(http.StatusForbidden)
	recorded.BodyIs("SSL Required.")
}

func TestTrustXFPHeaderForceSSLMiddleware(t *testing.T) {
	handler := NewAPI(&Middleware{
		TrustXFPHeader: true,
	})

	getRequest := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	getRequest.Header.Set("X-Forwarded-Proto", "http")
	recordedGet := test.RunRequest(t, handler, getRequest)

	recordedGet.CodeIs(http.StatusForbidden)
	recordedGet.BodyIs("SSL Required.")

	postRequest := test.MakeSimpleRequest("POST", "http://localhost/", simplePostData)
	postRequest.Header.Set("X-Forwarded-Proto", "http")
	recordedPost := test.RunRequest(t, handler, postRequest)

	recordedPost.CodeIs(http.StatusForbidden)
	recordedPost.BodyIs("SSL Required.")
}

func TestGetEnable301RedirectsForceSSLMiddleware(t *testing.T) {
	handler := NewAPI(&Middleware{
		Enable301Redirects: true,
	})

	getRequest := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	getRequest.Header.Set("X-Forwarded-Proto", "http")
	recordedGet := test.RunRequest(t, handler, getRequest)

	recordedGet.CodeIs(http.StatusMovedPermanently)

	postRequest := test.MakeSimpleRequest("POST", "http://localhost/", simplePostData)
	postRequest.Header.Set("X-Forwarded-Proto", "http")
	recordedPost := test.RunRequest(t, handler, postRequest)

	recordedPost.CodeIs(http.StatusMovedPermanently)
}

func TestMessageForceSSLMiddleware(t *testing.T) {
	message := "Custom message!"

	handler := NewAPI(&Middleware{
		Message: message,
	})

	getRequest := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recordedGet := test.RunRequest(t, handler, getRequest)

	recordedGet.CodeIs(http.StatusForbidden)
	recordedGet.BodyIs(message)

	postRequest := test.MakeSimpleRequest("POST", "http://localhost/", simplePostData)
	recordedPost := test.RunRequest(t, handler, postRequest)

	recordedPost.CodeIs(http.StatusForbidden)
	recordedPost.BodyIs(message)
}

func TestValidGetHTTPSRequestForceSSLMiddleware(t *testing.T) {
	handler := NewAPI(&Middleware{})

	getRequest := test.MakeSimpleRequest("GET", "https://localhost/", nil)
	recordedGet := test.RunRequest(t, handler, getRequest)

	recordedGet.CodeIs(http.StatusOK)
	recordedGet.BodyIs(`{"email":"edward@example.com","password":"password"}`)

	postRequest := test.MakeSimpleRequest("POST", "https://localhost/", simplePostData)
	recordedPost := test.RunRequest(t, handler, postRequest)

	recordedPost.CodeIs(http.StatusOK)
	recordedPost.BodyIs(`{"email":"edward@example.com","password":"password"}`)
}
