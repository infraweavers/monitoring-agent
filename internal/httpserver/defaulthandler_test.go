package httpserver

import (
	"net/http"
    "net/http/httptest"
    "io/ioutil"
    "testing"
)

func TestDefaultHandler(t *testing.T) {
    t.Run("with valid credentials, returns API root endpoints", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/", nil)
        request.SetBasicAuth("test", "secret")
        responseWriter := httptest.NewRecorder()

        DefaultHandler(responseWriter, request)
        response := responseWriter.Result()

        expectedStatusCode := 200
        outputStatusCode := response.StatusCode
        
        if expectedStatusCode != outputStatusCode {
            t.Error("Test Failed: Expected: {}, Got: {}", expectedStatusCode, outputStatusCode)
        }
        
        expectedBody := `{"endpoints": "/v1/"}`
        outputBodyBytes, _ := ioutil.ReadAll(response.Body)
        outputBody := string(outputBodyBytes)
        
        if expectedBody != outputBody {
            t.Error("Test Failed: Expected: {}, Got: {}", expectedBody, outputBody)
        }
    })
    
    t.Run("with no credentials, returns 401 unauthorized", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/", nil)
        responseWriter := httptest.NewRecorder()
        
        DefaultHandler(responseWriter, request)
        response := responseWriter.Result()

        expectedStatusCode := 401
        outputStatusCode := response.StatusCode
        
        if expectedStatusCode != outputStatusCode {
            t.Error("Test Failed: Expected: {}, Got: {}", expectedStatusCode, outputStatusCode)
        }
    })
    
    t.Run("with invalid credentials, returns 403 forbidden", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/", nil)
        request.SetBasicAuth("foo", "bar")
        responseWriter := httptest.NewRecorder()
        
        DefaultHandler(responseWriter, request)
        response := responseWriter.Result()

        expectedStatusCode := 403
        outputStatusCode := response.StatusCode
        
        if expectedStatusCode != outputStatusCode {
            t.Error("Test Failed: Expected: {}, Got: {}", expectedStatusCode, outputStatusCode)
        }
    })
}