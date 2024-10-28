package utils

import (
	"ads/internal/models"
	"ads/pkg/logger"

	"github.com/gin-gonic/gin"
)

func GinResponse(c *gin.Context, data models.Response) {

	switch data.Status {
	case 0:
		c.JSON(200, data.Data)
		return

	case 200:
		c.JSON(200, data.Data)
		return

	case 201:
		c.JSON(201, data.Data)
		return

	case 400:
		logger.Log.Println(data.Error.Error())
		c.JSON(400, models.InvalidInput)
		return

	case 402:
		logger.Log.Println(data.Error.Error())
		c.JSON(402, models.PaymentRequired)
		return

	case 404:
		logger.Log.Println(data.Error.Error())
		c.JSON(404, models.NotFound)
		return

	case 409:
		logger.Log.Println(data.Error.Error())
		c.JSON(409, models.Conflict)
		return

	case 500:
		logger.Log.Errorln(data.Error.Error())
		c.JSON(500, models.InternalServerError)
		return

	default:
		logger.Log.Errorln(data.Error.Error())
		c.JSON(500, models.InternalServerError)
		return
	}

}
