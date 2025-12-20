package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"whisper-api/mock"
	"whisper-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupHealthRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", handler.Health)
	return r
}

func testHealth(repo *repository.AuthRepository) (time.Duration, *httptest.ResponseRecorder) {
	handler := NewAuthHandler(repo)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(w, req)
	elapsed := time.Since(start)

	return elapsed, w
}

func TestHealth_DependenciesUp(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	elapsed, w := testHealth(repo)
	assert.Less(t, elapsed, 3*time.Second)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealth_DependenciesDown(t *testing.T) {
	mongoCli, _ := mongo.Connect(
		nil,
		options.Client().ApplyURI("mongodb://localhost:27099"),
	)

	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6399",
	})

	repo := repository.NewAuthRepository(
		mongoCli,
		redisCli,
		"testdb",
		"testcol",
	)

	elapsed, w := testHealth(repo)
	assert.Less(t, elapsed, 3*time.Second)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
