package controller

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	common.Response
	FavoriteList []FavoriteResponse `json:"video_list"`
}
type FavoriteResponse struct {
	Id            uint                  `json:"id"`
	User          UserInfoQueryResponse `json:"author"`
	PlayUrl       string                `json:"play_url"`
	CoverUrl      string                `json:"cover_url"`
	FavoriteCount uint                  `json:"favorite_count"`
	CommentCount  uint                  `json:"comment_count"`
	IsFavorite    bool                  `json:"is_favorite"`
	Title         string                `json:"title"`
}

func Favorite(c *gin.Context) {
	// 第一次使用先创建数据库
	dao.SqlSession.AutoMigrate(&model.Video{})
	dao.SqlSession.AutoMigrate(&model.Favorite{})

	// 获取videoId和actionType
	videoId := c.Query("video_id")
	actionType, _ := strconv.Atoi(c.Query("action_type"))

	// 从session中获取user
	rawUser := service.GetCurrentUser(c)

	// 判断行为,1是点赞,2是取消点赞
	if actionType == 1 {
		// 修改数据 点赞数+1
		err := service.FavoriteAction(videoId, rawUser.ID, 1)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				common.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			)
			return
		}
	} else if actionType == 2 {
		// 修改数据 点赞数-1
		err := service.FavoriteAction(videoId, rawUser.ID, -1)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				common.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			)
			return
		}
	} else {
		// 返回异常响应信息
		// 没有相关模式
		c.JSON(http.StatusBadRequest, CommentLeavingResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "请正确选择模式",
			},
		})
		return
	}

	// 返回响应结果
	c.JSON(
		http.StatusOK,
		common.Response{
			StatusCode: 0,
			StatusMsg:  "操作成功",
		},
	)

}

func FavoriteList(c *gin.Context) {
	// 第一次使用先创建数据库
	dao.SqlSession.AutoMigrate(&model.Video{})
	dao.SqlSession.AutoMigrate(&model.Favorite{})

	// 获取user_id
	userId := c.Query("user_id")

	// 根据 userId 获取视频信息
	videos, err := service.GetVideoByUserId(userId)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	var favorites []FavoriteResponse
	for _, video := range videos {
		// 根据author_id获取作者信息
		user, err := UserInfoService(strconv.Itoa(int(video.AuthorId)))
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		// 填充视频列表
		var favorite = FavoriteResponse{
			Id:            video.ID,
			User:          user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
		favorites = append(favorites, favorite)
	}

	// 返回响应结果
	c.JSON(
		http.StatusOK,
		FavoriteListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "操作成功",
			},
			FavoriteList: favorites,
		},
	)

}
