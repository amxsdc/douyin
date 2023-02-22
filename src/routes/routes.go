package routes

import (
	"douyin/src/controller"
	"douyin/src/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 设置生成session的密钥
	store := cookie.NewStore([]byte("douyin"))

	r.Use(sessions.Sessions("SESSIONID", store))

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
			commentGroup.POST("/action", middleware.JwtMiddleware(), controller.Comment)    // 发表｜删除评论
			commentGroup.GET("/list/", middleware.JwtMiddleware(), controller.CommentsList) // 评论列表
		}

		favoriteGroup := douyinGroup.Group("/favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Favorite)  // 点赞
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.FavoriteList) // 喜欢列表
		}

	}

	return r
}
