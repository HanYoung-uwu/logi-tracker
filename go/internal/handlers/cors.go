package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSPreflightHandler(c *gin.Context) {
	origin := c.GetHeader("Origin")

	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Vary", "Origin")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "content-type")
	c.JSON(http.StatusOK, "success")
}

func AllowCORSHandler(c *gin.Context) {
	origin := c.GetHeader("Origin")

	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Vary", "Origin")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}