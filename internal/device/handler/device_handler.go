package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/nms-go/internal/device/service"
)

type DeviceHandler struct {
	service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var req service.RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	device, err := h.service.RegisterDevice(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, device)
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	devices, total, err := h.service.ListDevices(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":       devices,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	device, err := h.service.GetDevice(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "device not found"})
		return
	}

	c.JSON(200, device)
}
