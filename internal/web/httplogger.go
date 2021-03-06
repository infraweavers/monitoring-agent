package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"monitoringagent/internal/logwrapper"
	"net/http"
	"net/http/httptest"
)

func httpRequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		bodyBytes, err := ioutil.ReadAll(request.Body)
		request.Body.Close()
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf(
				"%d Failed to read body from HTTP request: %#v",
				http.StatusInternalServerError,
				err.Error(),
			)))
			return
		}

		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		handler.ServeHTTP(responseWriter, request)

		logwrapper.LogHTTPRequest(
			request.RemoteAddr,
			request.Host,
			request.Method,
			request.URL.String(),
			request.Header,
			request.Proto,
			request.ContentLength,
			string(bodyBytes),
		)
	})
}

func httpResponseLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		responseRecorder := httptest.NewRecorder()
		handler.ServeHTTP(responseRecorder, request)

		response := responseRecorder.Result()
		bodyBytes, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf(
				"%d Failed to read body from HTTP response recorder: %#v",
				http.StatusInternalServerError,
				err.Error(),
			)))
			return
		}

		logwrapper.LogHTTPResponse(
			response.Status,
			response.Header,
			response.Proto,
			string(bodyBytes),
		)

		for key, value := range response.Header {
			responseWriter.Header()[key] = value
		}
		responseWriter.WriteHeader(responseRecorder.Code)
		responseRecorder.Body.WriteTo(responseWriter)
	})
}
