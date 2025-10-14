package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnregisterUser(t *testing.T) {
	os.Setenv("MONGO_CONNECTION", "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.3")
	os.Setenv("WHISPER_ADMIN_TOKEN", "admin123")
  	os.Setenv("WHISPER_DB", "whisper_test")

  	router := SetupRouter()
	t.Run("Non exisiting user", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/unregister/404", bytes.NewBufferString(""))
      req.Header.Set("X-Admin-Token", "admin123")
      w := httptest.NewRecorder()
      router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)	
	})

	t.Run("Successful unregistertion", func(t *testing.T) {
		os.Setenv("WHISPER_KEY", "testkey123")

      body := `{"owner":"` + testPhone + `","subject":"test","subscribers":["a","b"]}`
      req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
      req.Header.Set("X-Admin-Token", "admin123")
      w := httptest.NewRecorder()
      router.ServeHTTP(w, req)

      assert.Equal(t, http.StatusOK, w.Code)
      assert.NotEmpty(t, w.Body.String())

		unregister_req := httptest.NewRequest("POST", fmt.Sprintf("/unregister/%s", w.Body.String()), bytes.NewBufferString(""))
      unregister_req.Header.Set("X-Admin-Token", "admin123")
      unregistertion_w := httptest.NewRecorder()
      router.ServeHTTP(unregistertion_w, unregister_req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
