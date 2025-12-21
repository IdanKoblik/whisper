package endpoints

import (
	"whisper-api/models"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register godoc
// @Summary     Register new API token
// @Description Create a new API token for authentication. Requires admin privileges. Returns the generated API token.
// @Tags        admin
// @Accept      json
// @Produce     plain
// @Security    AdminAuth
// @Success     201 {string} string "API token created successfully"
// @Failure     400 {string} string "Failed to create token"
// @Failure     401 {string} string "Unauthorized - Admin access required"
// @Router      /admin/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	apiToken := uuid.NewString()
	data := &models.AuthModel{
		ApiToken: utils.HashToken(apiToken),
		Devices:  []string{},
	}

	err := h.Repo.CreateToken(c.Request.Context(), data)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(201, apiToken)
}
