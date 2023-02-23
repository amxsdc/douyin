package controller

import (
	"douyin/src/common"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FollowingResponse struct {
	StatusCode string             `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	UserList   []service.UserAttr `json:"user_list"`
}

type FollowersResponse struct {
	StatusCode string             `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	UserList   []service.UserAttr `json:"user_list"`
}

func FollowingList(c *gin.Context) {
	UserIdStr := c.Query("user_id")
	//获取关注列表信息
	FollowingList, err := FollowingService(UserIdStr)
	if err != nil {
		c.JSON(http.StatusOK, FollowingResponse{
			StatusCode: "1",
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FollowingResponse{
		StatusCode: "0",
		StatusMsg:  "success",
		UserList:   FollowingList,
	})

}

// FollowingService 获取关注列表主函数
func FollowingService(UserIdStr string) ([]service.UserAttr, error) {
	//验证id合法性
	if UserIdStr == "" {
		return nil, fmt.Errorf("invaild userid")
	}
	//格式转换
	UserId, err := strconv.ParseUint(UserIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	FollowingList, err := service.ObtainFollowingList(uint(UserId))
	if err != nil {
		return nil, err
	}
	return FollowingList, nil
}

func FollowersList(c *gin.Context) {

	UserIdStr := c.Query("user_id")
	//获取关注列表信息
	FollowersList, err := FollowerService(UserIdStr)
	if err != nil {
		c.JSON(http.StatusOK, FollowersResponse{
			StatusCode: "1",
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FollowersResponse{
		StatusCode: "0",
		StatusMsg:  "success",
		UserList:   FollowersList,
	})

}

// FollowerService 获取粉丝列表主函数
func FollowerService(UserIdStr string) ([]service.UserAttr, error) {
	//验证id合法性
	if UserIdStr == "" {
		return nil, fmt.Errorf("invaild userid")
	}
	//格式转换
	UserId, err := strconv.ParseUint(UserIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	FollowersList, err := service.ObtainFollowerList(uint(UserId))
	if err != nil {
		return nil, err
	}
	return FollowersList, nil
}

// FollowAction 关注操作
func FollowAction(c *gin.Context) {
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")
	//获取当前用户id
	UserData := service.GetCurrentUser(c)

	err := FollowActionService(toUserIdStr, actionTypeStr, UserData.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})

}

// FollowActionService 关注操作主函数
func FollowActionService(toUserIdStr string, actionTypeStr string, userId uint) error {
	//验证合法性
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 64)
	if err != nil {
		return err
	}
	actionType, err := strconv.ParseUint(actionTypeStr, 10, 64)
	if err != nil {
		return err
	}
	if actionType != 1 && actionType != 2 {
		return fmt.Errorf("invaild actionType")
	}

	//关注操作
	err = service.FollowAction(userId, uint(toUserId), uint(actionType))
	return err
}
