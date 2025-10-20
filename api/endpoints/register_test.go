package endpoints

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"whisper-api/db"
	"whisper-api/mock"
	"whisper-api/services"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const testPhone = "+1 2345678901"

func cleanupUser(t *testing.T, coll *mongo.Collection, owner string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.DeleteOne(ctx, map[string]string{"_id": owner})
	if err != nil {
		t.Fatalf("failed to cleanup user %s: %v", owner, err)
	}
}

func TestRegisterEndpoint_Handle(t *testing.T) {
	cfg := mock.ConfigMock(t)
	client, err := db.MongoConnection(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	collection := client.Database(cfg.Mongo.Database).Collection("users")
	cleanupUser(t, collection, testPhone)
	router := SetupRouter(&cfg)

	t.Run("Invalid body", func(t *testing.T) {
		body := `{}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("X-Admin-Token", "admin123")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)	
	})
	
	t.Run("Unauthorized admin token", func(t *testing.T) {
		body := `{"owner":"` + testPhone + `","subject":"test","subscribers":["a","b"]}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("X-Admin-Token", "wrongtoken")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("Invalid phone number", func(t *testing.T) {
		body := `{"owner":"12345","subject":"test","subscribers":["a","b"]}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("X-Admin-Token", "admin123")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid phone number")
	})

	t.Run("Successful registration", func(t *testing.T) {
		os.Setenv("WHISPER_KEY", "testkey123")

		body := `{"owner":"` + testPhone + `"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("X-Admin-Token", cfg.AdminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println(w.Body.String())
		assert.NotEmpty(t, w.Body.String())

		var user services.User
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var resp RegisterResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		tokenHash := sha256.Sum256([]byte(resp.Token))
      tokenStr := hex.EncodeToString(tokenHash[:])
		err = collection.FindOne(ctx, map[string]string{"_id": tokenStr}).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, testPhone, user.Owner)

		cleanupUser(t, collection, user.Owner)
	})

}
