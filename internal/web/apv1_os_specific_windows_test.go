package web

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSSpecificApiHandler(t *testing.T) {

	TestSetup()
	defer TestTeardown()

	t.Run("When supplied an executable, returns HTTP status 200 and expected script output", func(t *testing.T) {
		testRequest := "{ \"CounterPath\": \"\\\\Memory\\\\Available MBytes\" }"

		output := RunTestRequest(t, http.MethodPost, "/v1/os_specific", strings.NewReader(testRequest))

		assert := assert.New(t)
		assert.Equal(http.StatusOK, output.ResponseStatus, "Response code should be OK")
		assert.NotNil(output.ResponseBody)
	})
}
