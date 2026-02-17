package apigateway

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/config_mgt"
	"github.com/yourorg/nms-go/internal/device/handler"
	"github.com/yourorg/nms-go/internal/device/repository"
	"github.com/yourorg/nms-go/internal/device/service"
	"github.com/yourorg/nms-go/internal/features/execution"
	"github.com/yourorg/nms-go/internal/features/monitoring"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB, monitoringHandler *monitoring.Handler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize dependencies
	deviceRepo := repository.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceService)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		devices := v1.Group("/devices")
		{
			devices.GET("", deviceHandler.ListDevices)
			devices.POST("", deviceHandler.RegisterDevice)
			devices.GET("/:id", deviceHandler.GetDevice)
		}

		// Config Management routes
		sshAdapter := config_mgt.NewSSHAdapter()
		configService := config_mgt.NewConfigService(deviceService, sshAdapter)
		configHandler := config_mgt.NewConfigHandler(configService)

		configGroup := v1.Group("/config")
		{
			configGroup.POST("/execute", configHandler.ExecuteCommand)
		}

		// Execution feature (Realtime)
		execService := execution.NewExecutionService()
		execHandler := execution.NewExecutionHandler(execService)
		v1.POST("/realtime/execute", execHandler.ExecuteCommand)
		v1.POST("/realtime/stats", execHandler.GetStats)

		// Monitoring feature (Background)
		v1.POST("/inventory/sync", monitoringHandler.SyncInventory)
	}

	return r
}
