package httpserver

import (
	"io/ioutil"
	"mama/internal/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServer *httptest.Server

type TestCredential struct {
	Username string
	Password string
}

// TestHTTPResponse is a struct consisting of responseBody and responseStatus
type TestHTTPResponse struct {
	ResponseStatus int
	ResponseBody   string
}

// TestSetup instantiates a router and test HTTP server
func TestSetup() {
	configuration.TestingInitialise()
	router := NewRouter()
	testServer = httptest.NewServer(router)
}

// GetServerURL returns the hostname and port for a running test server
func GetTestServerURL(t *testing.T) string {
	if testServer == nil {
		t.Fatal("HTTP test server URL requested via GetServerURL before Setup() called")
	}

	return testServer.URL
}

// TestTeardown closes the http server
func TestTeardown() {
	testServer.Close()
}

func BuildTestHTTPRequest(t *testing.T, method string, url string) *http.Request {
	request, error := http.NewRequest(method, GetTestServerURL(t)+url, nil)
	if error != nil {
		t.Fatalf(error.Error())
	}
	return request
}

func TestHTTPRequestWithCredentials(t *testing.T, request *http.Request, username string, password string) TestHTTPResponse {
	request.SetBasicAuth(username, password)
	return TestHTTPRequest(t, request)
}

func TestHTTPRequestWithDefaultCredentials(t *testing.T, request *http.Request) TestHTTPResponse {
	request.SetBasicAuth("test", "secret")
	return TestHTTPRequest(t, request)
}

func TestHTTPRequest(t *testing.T, request *http.Request) TestHTTPResponse {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	return TestHTTPResponse{
		ResponseStatus: response.StatusCode,
		ResponseBody:   string(body),
	}
}
