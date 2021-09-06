package web

import (
	"io/ioutil"
	"monitoringagent/internal/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServer *httptest.Server

// StdIn is a struct for use in setting up a test case that passes standard input (StdIn) to an endpoint
type StdIn struct {
	StdIn string
}

// Timeout is a struct for use in setting up a test case that passes standard input (StdIn) to an endpoint
type Timeout struct {
	Timeout string
}

// ExpectedResult is a struct for use in setting up a test case that defines the output string that an endpoint should return
type ExpectedResult struct {
	Output string
}

// StdInSignature is a struct for storing the signature of the script passed to stdin
type StdInSignature struct {
	StdInSignature string
}

// ScriptToRun is a struct that defines a script (path and arguments) to be passed to an endpoint under test
type ScriptToRun struct {
	Path string
	Args []string
}

// TestCredential is a struct for use in setting up test cases
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

// GetTestServerURL returns the hostname and port for a running test server
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

// BuildTestHTTPRequest creates a request for testing purposes. Automatically prepends the test server url.
func BuildTestHTTPRequest(t *testing.T, method string, url string) *http.Request {
	request, error := http.NewRequest(method, GetTestServerURL(t)+url, nil)
	if error != nil {
		t.Fatalf(error.Error())
	}
	return request
}

// TestHTTPRequestWithCredentials executes an HTTP request, with provided basic auth credentials
func TestHTTPRequestWithCredentials(t *testing.T, request *http.Request, username string, password string) TestHTTPResponse {
	request.SetBasicAuth(username, password)
	return TestHTTPRequest(t, request)
}

// TestHTTPRequestWithDefaultCredentials executes an HTTP request with the default credentials specified in the configuration file
func TestHTTPRequestWithDefaultCredentials(t *testing.T, request *http.Request) TestHTTPResponse {
	request.SetBasicAuth(configuration.Settings.Authentication.Username, configuration.Settings.Authentication.Password)
	return TestHTTPRequest(t, request)
}

// TestHTTPRequest executes an HTTP request with the provided request. Returns an HTTP Response.
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
