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

// CommentLeavingResponse 外部controller返回的响应信息
type CommentLeavingResponse struct {
	common.Response
	CommentResponse CommentResponse `json:"comment"`
}

// CommentResponse CommentLeavingResponse里面的Comment
type CommentResponse struct {
	Id         uint                  `json:"id"`
	UserInfo   UserInfoQueryResponse `json:"user"`
	Content    string                `json:"content"`
	CreateDate string                `json:"create_date"`
}

type CommentListResponse struct {
	common.Response
	Comments []CommentResponse `json:"comment_list"`
}

func Comment(c *gin.Context) {
	// 第一次使用先创建数据库
	dao.SqlSession.AutoMigrate(&model.Comment{})

	// 获取videoId, action_type等基本信息
	videoId := c.Query("video_id")
	actionType, _ := strconv.Atoi(c.Query("action_type"))

	// 从session中获取userId
	userId := strconv.Itoa(int(service.GetCurrentUser(c).ID))

	// 1为发布评论, 2为删除评论
	if actionType == 1 {
		commentText := c.Query("comment_text")
		comment, err := CommentLeaving(videoId, userId, commentText)
		if err != nil {
			c.JSON(http.StatusBadRequest, CommentLeavingResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		// 返回响应信息
		c.JSON(http.StatusOK, CommentLeavingResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "操作成功",
			},
			CommentResponse: comment,
		})
	} else if actionType == 2 {
		commentId := c.Query("comment_id")
		comment, err := CommentDeleting(videoId, userId, commentId)
		if err != nil {
			c.JSON(http.StatusBadRequest, CommentLeavingResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		// 返回响应信息
		c.JSON(http.StatusOK, CommentLeavingResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "操作成功",
			},
			CommentResponse: comment,
		})
	} else {
		// 返回异常响应信息
		c.JSON(http.StatusBadRequest, CommentLeavingResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "请正确选择模式",
			},
		})
		return
	}
}

func CommentLeaving(videoId string, userId string, commentText string) (CommentResponse, error) {
	// 查询用户信息
	user, err := UserInfoService(userId)

	// 添加评论
	commentId, createDate, err := service.AddComment(commentText, userId, videoId)

	// 返回响应消息
	return CommentResponse{
		Id:         commentId,
		UserInfo:   user,
		Content:    commentText,
		CreateDate: createDate,
	}, err
}

func CommentDeleting(videoId string, userId string, commentId string) (CommentResponse, error) {
	// 查询用户信息
	user, err := UserInfoService(userId)
	if err != nil {
		return CommentResponse{}, err
	}

	// 删除评论
	commentText, createDate, err := service.DeleteCommentById(commentId, userId, videoId)

	if err != nil {
		return CommentResponse{}, err
	}

	// 对commentId 进行格式转化
	comId, _ := strconv.Atoi(commentId)

	// 返回响应消息
	return CommentResponse{
		Id:         uint(comId),
		UserInfo:   user,
		Content:    commentText,
		CreateDate: createDate,
	}, err
}

func CommentsList(c *gin.Context) {
	// 第一次使用先创建数据库
	dao.SqlSession.AutoMigrate(&model.Comment{})

	// 获取videoId
	videoId := c.Query("video_id")

	// 根据videoId查询评论
	comments := service.ListAllComments(videoId)

	// 将Comment插入到CommentListResponse中
	var commentResponses []CommentResponse
	for _, comment := range comments {
		userId := comment.UserId
		userInfo, err := UserInfoService(strconv.Itoa(int(userId)))
		if err != nil {
			return
		}
		commentResponses = append(commentResponses, CommentResponse{
			Id:         comment.ID,
			UserInfo:   userInfo,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(
		http.StatusOK,
		CommentListResponse{
			Response: common.Response{
				StatusCode: http.StatusOK,
				StatusMsg:  "查询成功",
			},
			Comments: commentResponses,
		},
	)
}
