package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"server/internal/appcontext"
	"time"
)

func SetupRouter(appContext *appcontext.AppContext) *gin.Engine {
	r := gin.Default()

	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	api := r.Group("/api")
	{
		api.POST("/package", func(c *gin.Context) { HandlePackageRequest(c, appContext) })
		api.GET("/packs", func(c *gin.Context) { HandleGetPacksRequest(c, appContext) })
		api.POST("/packs", func(c *gin.Context) { HandlePacksSyncRequest(c, appContext) })
	}

	return r
}
