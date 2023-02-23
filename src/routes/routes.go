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
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo)
			userGroup.POST("/login/", controller.UserLogin)
			userGroup.POST("/register/", controller.UserRegister)
		}

		// publish路由组
		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action/", middleware.JwtMiddleware(), controller.Publish)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), controller.PublishList)

		}

		// feed
		douyinGroup.GET("/feed/", controller.Feed)

		favoriteGroup := douyinGroup.Group("favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Favorite)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.FavoriteList)
		}

		// comment路由组
		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/action/", middleware.JwtMiddleware(), controller.CommentAction)
			commentGroup.GET("/list/", middleware.JwtMiddleware(), controller.CommentList)
		}

		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)
			relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
			relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
		}

		messageGroup := douyinGroup.Group("/message")
		{
			messageGroup.POST("/action/", middleware.JwtMiddleware(), controller.SendingMessage) // 发送消息
			messageGroup.GET("/chat/", middleware.JwtMiddleware(), controller.MessageLog)        // 聊天记录
		}

	}

	return r
}
