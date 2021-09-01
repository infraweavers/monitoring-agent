package web

import (
	"fmt"
	"io/ioutil"
	"monitoringagent/internal/configuration"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var routePathTemplates []string
var router *mux.Router
var server *httptest.Server

func setup(t *testing.T) {
	router = NewRouter()
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			routePathTemplates = append(routePathTemplates, pathTemplate)
		}

		return nil
	})

	assert.Nil(t, err)

	server = httptest.NewServer(router)
	configuration.TestingInitialise()
}

func teardown() {
	server.Close()
}

func testRoute(t *testing.T, path string) {
	t.Run(fmt.Sprintf("returns HTTP status 200 for path %s", path), func(t *testing.T) {
		assert := assert.New(t)
		request, err := http.NewRequest(http.MethodGet, path, nil)
		assert.Nil(err)

		request.SetBasicAuth("test", "secret")

		response, err := http.DefaultClient.Do(request)
		assert.Nil(err)

		defer response.Body.Close()

		assert.Equal(http.StatusOK, response.StatusCode)

		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(err)

	})
}

func TestAllRoutes(t *testing.T) {
	setup(t)

	for _, path := range routePathTemplates {
		testRoute(t, server.URL+path)
	}

	teardown()
}
