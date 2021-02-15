package httpserver

import (
	"fmt"
	"io/ioutil"
	"mama/internal/configuration"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

	if err != nil {
		t.Fatal(err)
	}

	server = httptest.NewServer(router)
	configuration.TestingInitialise()
}

func teardown() {
	server.Close()
}

func testRoute(t *testing.T, path string) {
	t.Run(fmt.Sprintf("returns HTTP status 200 for path %s", path), func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "secret")
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Received non-%d response: %d\n", http.StatusOK, response.StatusCode)
		}

		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

	})
}

func TestAllRoutes(t *testing.T) {
	setup(t)

	for _, path := range routePathTemplates {
		testRoute(t, server.URL+path)
	}

	teardown()
}
