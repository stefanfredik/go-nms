package olt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is the Gin HTTP handler for OLT API endpoints.
// All endpoints accept a JSON request body with SNMP target details —
// no device registry is required in go-nms.
type Handler struct {
	service OLTService
}

// NewHandler creates a new OLT HTTP handler.
func NewHandler(service OLTService) *Handler {
	return &Handler{service: service}
}

// GetSystemMetrics handles POST /api/v1/olt/system
//
// Returns system-level metrics (CPU, memory, uptime, temperature) for the OLT
// specified in the request body.
func (h *Handler) GetSystemMetrics(c *gin.Context) {
	var req GetSystemMetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	metrics, err := h.service.GetSystemMetrics(c.Request.Context(), req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetPONPorts handles POST /api/v1/olt/pon-ports
//
// Returns metrics for all PON ports on the OLT specified in the request body.
func (h *Handler) GetPONPorts(c *gin.Context) {
	var req GetPONPortsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	ports, err := h.service.GetPONPorts(c.Request.Context(), req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ports)
}

// GetONTs handles POST /api/v1/olt/onts
//
// Returns metrics for all ONTs on the OLT specified in the request body.
// Set pon_port > 0 in the body to filter by a specific PON port.
func (h *Handler) GetONTs(c *gin.Context) {
	var req GetONTsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if req.PONPort < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pon_port must be >= 0"})
		return
	}

	onts, err := h.service.GetONTs(c.Request.Context(), req.Target, req.PONPort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, onts)
}

// GetONTStatus handles POST /api/v1/olt/ont-status
//
// Returns ONTs categorized by their operational status (Up/Down).
func (h *Handler) GetONTStatus(c *gin.Context) {
	// usage: Use GetSystemMetricsRequest since it only contains Target, which is exactly what we need.
	var req GetSystemMetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	status, err := h.service.GetONTStatus(c.Request.Context(), req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// RegisterRoutes registers all OLT routes on the given Gin router group.
func RegisterRoutes(group *gin.RouterGroup, service OLTService) {
	h := NewHandler(service)

	oltGroup := group.Group("/olt")
	{
		// POST /api/v1/olt/system     — system metrics (CPU, memory, uptime, temperature)
		oltGroup.POST("/system", h.GetSystemMetrics)

		// POST /api/v1/olt/pon-ports  — PON port status and optical power
		oltGroup.POST("/pon-ports", h.GetPONPorts)

		// POST /api/v1/olt/onts       — ONT list (filter by pon_port in body)
		oltGroup.POST("/onts", h.GetONTs)

		// POST /api/v1/olt/ont-status — ONT status list (up/down)
		oltGroup.POST("/ont-status", h.GetONTStatus)
	}
}
