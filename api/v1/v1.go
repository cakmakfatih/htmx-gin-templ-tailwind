package api

import "github.com/gin-gonic/gin"

func RegisterApi(router *gin.Engine) {
	v1 := router.Group("/api/v1/")

	registerAuth(v1)
}
