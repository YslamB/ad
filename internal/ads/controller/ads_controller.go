package controller

import (
	"ads/internal/ads/usecase"
	"ads/internal/cache"
	"ads/internal/models"
	"ads/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdController struct {
	usecase *usecase.AdUsecase
}

func NewAdController(db *pgxpool.Pool, rClient *cache.RedisCache) *AdController {
	return &AdController{usecase.NewAdUsecase(db, rClient)}
}

func (ctrl *AdController) GetAdByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.MustGet("paramID").(int)
	data := ctrl.usecase.GetAdByID(&ctx, id)
	utils.GinResponse(c, data)
}

func (ctrl *AdController) GetAds(c *gin.Context) {
	ctx := c.Request.Context()
	page := c.MustGet("page").(int)
	limit := c.MustGet("limit").(int)
	orderBy := c.MustGet("order_by").(string)
	data := ctrl.usecase.GetAds(&ctx, page, limit, orderBy)
	utils.GinResponse(c, data)
}

func (ctrl *AdController) CreateAd(c *gin.Context) {
	ctx := c.Request.Context()
	var reqBody models.CreateAd
	validationError := c.BindJSON(&reqBody)

	if validationError != nil {
		utils.GinResponse(c, models.Response{Status: 400, Error: validationError})
		return
	}

	data := ctrl.usecase.CreateAd(&ctx, reqBody)
	utils.GinResponse(c, data)
}

func (ctrl *AdController) EditAd(c *gin.Context) {
	ctx := c.Request.Context()
	var reqBody models.CreateAd
	validationError := c.BindJSON(&reqBody)

	if validationError != nil {
		utils.GinResponse(c, models.Response{Status: 400, Error: validationError})
		return
	}

	data := ctrl.usecase.EditAd(&ctx, reqBody)
	utils.GinResponse(c, data)
}

func (ctrl *AdController) DeleteAdByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.MustGet("paramID").(int)
	data := ctrl.usecase.DeleteAdByID(&ctx, id)
	utils.GinResponse(c, data)
}
