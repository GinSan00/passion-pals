package authhttp

import (
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

func Register(router *gin.Engine, authService Auth) {
	router.POST("/login", authService.Login)
	router.POST("/register", authService.Register)
}
