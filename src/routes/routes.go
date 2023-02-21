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
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo) // 用户信息
			userGroup.POST("/login/", controller.UserLogin)                     // 用户登录
			userGroup.POST("/register/", controller.UserRegister)               // 用户注册
		}

		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/", middleware.JwtMiddleware(), controller.Comment)
		}

	}

	return r
}
