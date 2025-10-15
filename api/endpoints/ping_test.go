package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestPingEndpoint(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(&cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
