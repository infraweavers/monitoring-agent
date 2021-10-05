package web

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"monitoringagent/internal/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServer *httptest.Server

// TestCredential is a struct for representing a set of credentials (user name and password), used in setting up test cases
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

// RunTestRequest Runs the test JSON against the relative URL on the test HTTP server
func RunTestRequest(t *testing.T, method, url string, body io.Reader) TestHTTPResponse {
	request, _ := http.NewRequest(method, GetTestServerURL(t)+url, body)
	return TestHTTPRequestWithDefaultCredentials(t, request)
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

// JsonSerialize returns a byte buffer of an input object
func JsonSerialize(inputObject interface{}) *bytes.Buffer {
	byteArray, _ := json.Marshal(inputObject)
	byteArrayBuffer := bytes.NewBuffer(byteArray)
	return byteArrayBuffer
}
