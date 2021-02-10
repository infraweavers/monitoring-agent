package httpserver

import (
	"net/http"
    "net/http/httptest"
    "testing"
)

func TestDefaultHandler(t *testing.T) {
    t.Run("returns API root endpoints", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/", nil)
        responseWriter := httptest.NewRecorder()
        request.SetBasicAuth("test", "secret")

        DefaultHandler(responseWriter, request)

        expected := `{"endpoints": "/v1/"}`
        output := responseWriter.Body.String()

        if expected != output {
            t.Error("Test Failed: Expected: {}, Got: {}", expected, output)
        }
    })
}