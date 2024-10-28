package delivery

import (
	"ads/internal/ads/controller"
	"ads/internal/cache"
	"ads/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterAdsRoutes(router *gin.Engine, dbPool *pgxpool.Pool, rClient *cache.RedisCache) {

	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup AD Routes
	rAd := router.Group("/ads")
	adController := controller.NewAdController(dbPool, rClient)
	rAd.GET("", middleware.PageLimitOrderSet, adController.GetAds)
	rAd.GET("/:id", middleware.ParamIDToInt, adController.GetAdByID)
	rAd.PUT("/update", adController.EditAd)
	rAd.POST("/create", adController.CreateAd)
	rAd.DELETE("/delete/:id", middleware.ParamIDToInt, adController.DeleteAdByID)

}
