package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
)

func AddMessage(toUserId uint, fromUserId uint, content string) error {

	var message = model.Message{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
	}

	if err := dao.SqlSession.Model(&model.Message{}).Create(&message).Error; err != nil {
		return common.ErrorInsertion
	}

	return nil
}

func CheckMessage(toUserId string, fromUserId uint) ([]model.Message, error) {
	var messages []model.Message
	if err := dao.SqlSession.Model(&model.Message{}).Where("to_user_id=? and from_user_id=?", toUserId, fromUserId).Find(&messages).Error; err != nil {
		return []model.Message{}, common.ErrorSelection
	}
	return messages, nil
}
