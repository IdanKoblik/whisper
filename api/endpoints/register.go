package endpoints

import (
	"os"
	"regexp"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterEndpoint struct {
	service *services.UserService	
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
// @Success 200 {string} string "JWT Token"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /register [post]
func (endpoint RegisterEndpoint) Handle(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	if adminToken != os.Getenv("WHISPER_ADMIN_TOKEN") {
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
	key := os.Getenv("WHISPER_KEY")

	signedToken, err := token.SignedString([]byte(key))
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

	c.String(200, user.Token)
}
