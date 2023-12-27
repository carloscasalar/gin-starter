package status_test

import (
	"encoding/json"
	"fmt"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/status"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	t.Run("status should be 200", func(t *testing.T) {
		ts := httptest.NewServer(setupServer())
		defer ts.Close()

		resp, err := http.Get(fmt.Sprintf("%s/status", ts.URL))

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("body should contain a message telling the service is up and running", func(t *testing.T) {
		ts := httptest.NewServer(setupServer())
		defer ts.Close()

		resp, err := http.Get(fmt.Sprintf("%s/status", ts.URL))

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		payload := getPayload(t, resp)
		assert.Equal(t, "server is ready and healthy", payload.Message)
	})
}

func getPayload(t *testing.T, resp *http.Response) statusResponse {
	payload := statusResponse{}
	err := json.NewDecoder(resp.Body).Decode(&payload)
	require.NoError(t, err)
	return payload
}

func setupServer() http.Handler {
	r := gin.Default()

	r.GET("/status", status.Handler)
	return r
}

type statusResponse struct {
	Message string `json:"message"`
}