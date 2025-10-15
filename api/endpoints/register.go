package endpoints

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"
	"whisper-api/config"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterEndpoint struct {
	service *services.UserService
	cfg *config.Config
}

type RegisterResponse struct {
	Token string `json:"token"`
	Key string `json:"key"` 
}

const PATTERN = `^\+\d{1,3} \d{7,12}$`

// RegisterEndpoint godoc
// @Summary Register a new user
// @Description Requires X-Admin-Token header and RawUser JSON body
// @Tags Users
// @Accept json
// @Produce json
// @Param X-Admin-Token header string true "Admin Token"
// @Param rawUser body services.RawUser true "User data"
// @Success 200 {string} string "JWT Token and Signature key"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /register [post]
func (endpoint RegisterEndpoint) Handle(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	if adminToken != endpoint.cfg.AdminToken {
		c.String(401, "Unauthorized: invalid admin token")
		return
	}

	var rawUser services.RawUser 
	err := c.ShouldBindBodyWithJSON(&rawUser)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	re := regexp.MustCompile(PATTERN)
	if !re.MatchString(rawUser.Owner) {
		c.String(400, "Invalid phone number")
		return
	}

	jwtData := jwt.MapClaims {
		"owner": rawUser.Owner,
		"subject": rawUser.Subject,
		"subscribers": rawUser.Subscribers,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	rawID := fmt.Sprintf("%s-%s", uuid.New(), time.Now().String())
	identifier := sha256.Sum256([]byte(rawID))
	
	signedToken, err := token.SignedString(identifier[:])
	if err != nil {
		c.String(400, err.Error())
		return
	}

	user := services.User {
		Owner: rawUser.Owner,
		Token: signedToken, 
		Subject: rawUser.Subject,
		Subscribers: rawUser.Subscribers,
	}

	err = endpoint.service.RegisterUser(&user)	
	if err != nil {
		c.String(400, err.Error())
		return
	}

	response := RegisterResponse {
		Token: user.Token,
		Key: hex.EncodeToString(identifier[:]),
	}

	c.JSON(200, response)
}
