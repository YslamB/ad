package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParamIDToInt(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.AbortWithStatus(400)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	c.Set("paramID", id)
	c.Next()
}

func PageLimitOrderSet(c *gin.Context) {
	orderBy := c.Query("order") // -time, time, -price, price
	switch orderBy {
	case "time":
		orderBy = "order by created_at"
	case "price":
		orderBy = "order by price"
	case "-price":
		orderBy = "order by price desc"
	default:
		orderBy = "order by created_at desc"
	}

	pageStr := c.Query("page")
	countStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(countStr)
	if err != nil {
		limit = 10
	}

	c.Set("page", page)
	c.Set("limit", limit)
	c.Set("order_by", orderBy)
	c.Next()
}

// func Cors(c *gin.Context) {
// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

// 	if c.Request.Method == "OPTIONS" {
// 		c.AbortWithStatus(204)
// 		return
// 	}

// 	c.Next()
// }
