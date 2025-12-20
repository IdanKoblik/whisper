package endpoints

import (
	"github.com/gin-gonic/gin"
)

type DeviceRequest struct {
	Device string `json:"device"`
}

// AddDevice godoc
// @Summary     Add device
// @Description Add a new device to the authenticated user's device list
// @Tags        devices
// @Accept      json
// @Produce     plain
// @Security    ApiKeyAuth
// @Param       request body DeviceRequest true "Device information"
// @Success     201 {string} string "Device added successfully"
// @Failure     400 {string} string "Invalid request or failed to add device"
// @Failure     401 {string} string "Unauthorized - Invalid or missing token"
// @Router      /api/devices [post]
func (h *AuthHandler) AddDevice(c *gin.Context) {
	var request DeviceRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = h.Repo.AddDeviceID(c.Request.Context(), c.GetString("token"), request.Device)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(201, "Added new device")
}

// RemoveDevice godoc
// @Summary     Remove device
// @Description Remove a device from the authenticated user's device list
// @Tags        devices
// @Accept      json
// @Produce     plain
// @Security    ApiKeyAuth
// @Param       request body DeviceRequest true "Device information"
// @Success     200 {string} string "Device removed successfully"
// @Failure     400 {string} string "Invalid request or failed to remove device"
// @Failure     401 {string} string "Unauthorized - Invalid or missing token"
// @Router      /api/devices [delete]
func (h *AuthHandler) RemoveDevice(c *gin.Context) {
	var request DeviceRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = h.Repo.RemoveDeviceID(c.Request.Context(), c.GetString("token"), request.Device)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(200, "Device removed successfully")
}

// GetDevice godoc
// @Summary     Get device
// @Description Check if a device exists and belongs to the authenticated user
// @Tags        devices
// @Accept      json
// @Produce     plain
// @Security    ApiKeyAuth
// @Param       id path string true "Device ID"
// @Success     200 {string} string "Device found"
// @Success     404 {string} string "Device not found"
// @Failure     400 {string} string "Invalid request or validation error"
// @Failure     401 {string} string "Unauthorized - Invalid or missing token"
// @Router      /api/devices/{id} [get]
func (h *AuthHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")
	found, err := h.Repo.ValidateDeviceID(c.Request.Context(), c.GetString("token"), deviceID)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	if found {
		c.String(200, "Device found")
	} else {
		c.String(404, "Device not found")
	}
}

// DeviceHandler handles device operations for backward compatibility
// This function routes to the appropriate handler based on HTTP method
func (h *AuthHandler) DeviceHandler(c *gin.Context) {
	method := c.Request.Method

	switch method {
	case "POST":
		h.AddDevice(c)
	case "DELETE":
		h.RemoveDevice(c)
	case "GET":
		h.GetDevice(c)
	}
}
