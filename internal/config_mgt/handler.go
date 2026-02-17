package config_mgt

import (
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service ConfigService
}

func NewConfigHandler(service ConfigService) *ConfigHandler {
	return &ConfigHandler{service: service}
}

type ExecuteCommandRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Command  string `json:"command" binding:"required"`
}

func (h *ConfigHandler) ExecuteCommand(c *gin.Context) {
	var req ExecuteCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// output is now interface{}, standard JSON marshaling will handle it (string or object)
	output, err := h.service.ExecuteCommand(c.Request.Context(), req.DeviceID, req.Command)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "output": output})
		return
	}

	c.JSON(200, gin.H{"output": output})
}
