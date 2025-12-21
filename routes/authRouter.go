package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/yaikob/goland-jwt/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("user/signup", controller.Signup())
	incomingRoutes.POST("user/login", controller.Login())
}
