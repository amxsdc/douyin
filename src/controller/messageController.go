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

type MessageLogResponse struct {
	common.Response
	MessageResponses []MessageResponse `json:"message_list"`
}
type MessageResponse struct {
	Id         uint   `json:"id"`
	ToUserId   uint   `json:"to_user_id"`
	FromUserId uint   `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

func SendingMessage(c *gin.Context) {
	// 第一次使用创建数据库
	dao.SqlSession.AutoMigrate(&model.Message{})

	// 获取toUserId
	toUserId, _ := strconv.Atoi(c.Query("to_user_id"))

	// 从session中获取fromUserId
	fromUserId := service.GetCurrentUser(c).ID

	// 获取actionType
	actionType, _ := strconv.Atoi(c.Query("action_type"))

	// 获取content
	content := c.Query("content")

	// 1-发送消息
	if actionType == 1 {
		// 存储消息内容
		err := service.AddMessage(uint(toUserId), fromUserId, content)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				common.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			)
		}
	} else {
		c.JSON(
			http.StatusBadRequest,
			common.Response{
				StatusCode: 1,
				StatusMsg:  "请正确选择模式",
			},
		)
		return
	}

	// 返回响应消息
	c.JSON(
		http.StatusOK,
		common.Response{
			StatusCode: 0,
			StatusMsg:  "发送成功",
		},
	)
}

func MessageLog(c *gin.Context) {
	// 第一次使用创建数据库
	dao.SqlSession.AutoMigrate(&model.Message{})

	// 获取toUserId
	toUserId := c.Query("to_user_id")

	// 从session中获取fromUserId
	fromUserId := service.GetCurrentUser(c).ID

	// 通过两个id进行查询
	messages, err := service.CheckMessage(toUserId, fromUserId)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.Response{
				StatusCode: 1,
				StatusMsg:  "没有消息",
			},
		)
		return
	}

	var messageResponses []MessageResponse
	// 封装消息到MessageLogResponse
	for _, message := range messages {
		messageResponse := MessageResponse{
			Id:         message.ID,
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		messageResponses = append(messageResponses, messageResponse)
	}

	// 返回响应
	c.JSON(
		http.StatusOK,
		MessageLogResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "获取成功",
			},
			MessageResponses: messageResponses,
		},
	)
}
