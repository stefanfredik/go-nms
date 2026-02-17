package monitoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *TargetStore
}

func NewHandler(store *TargetStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) SyncInventory(c *gin.Context) {
	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targets := make([]DeviceTarget, len(req.Targets))
	for i, t := range req.Targets {
		targets[i] = DeviceTarget{
			IP:       t.IP,
			Driver:   t.Driver,
			Username: t.Auth.Username,
			Password: t.Auth.Password,
			Port:     t.Auth.Port,
		}
	}

	h.store.ReplaceAll(targets)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(targets),
	})
}
