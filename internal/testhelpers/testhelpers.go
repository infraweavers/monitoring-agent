package testhelpers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var routePathTemplates []string
var router *mux.Router
var server *httptest.Server

// HTTPResponse is a struct consisting of responseBody and responseStatus
type HTTPResponse struct {
	ResponseStatus int
	ResponseBody   string
}

// Setup instantiates a router and test HTTP server
func Setup(handlerFunc http.HandlerFunc) {
	server = httptest.NewServer(handlerFunc)
}

// GetServerURL returns the hostname and port for a running test server
func GetServerURL(t *testing.T) string {
	if server == nil {
		t.Fatal("HTTP test server URL requested via GetServerURL before Setup() called")
	}

	return server.URL
}

// Teardown closes the http server
func Teardown() {
	server.Close()
}

// TestHTTPRequest executes an HTTP request and returns a struct containing the body and HTTP response code
func TestHTTPRequest(t *testing.T, request *http.Request) HTTPResponse {
	request.SetBasicAuth("test", "secret")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	return HTTPResponse{
		ResponseStatus: response.StatusCode,
		ResponseBody:   string(body),
	}
}
