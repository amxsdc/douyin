package routes

import (
	"douyin/src/controller"
	"douyin/src/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo)
			userGroup.POST("/login/", controller.UserLogin)
			userGroup.POST("/register/", controller.UserRegister)
		}

	}

	return r
}
