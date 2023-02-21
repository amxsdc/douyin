package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"strconv"
)

// AddComment 添加评论的持久化操作
func AddComment(commentText string, userId string, videoId string) (uint, string, error) {
	if len(commentText) == 0 {
		return 0, "", common.ErrorCommentNull
	}
	// 开启事务
	tx := dao.SqlSession.Begin()

	// 添加评论
	//将string类型数据转化为uint类型数据
	vId, _ := strconv.Atoi(videoId)
	uId, _ := strconv.Atoi(userId)
	comment := model.Comment{
		VideoId: uint(vId),
		UserId:  uint(uId),
		Content: commentText,
	}
	if err := dao.SqlSession.Model(&model.Comment{}).Create(&comment).Error; err != nil {
		// 评论添加失败
		tx.Rollback()
		return 0, "", common.ErrorInsertion
	}
	// 获取评论id
	var commentId []uint
	dao.SqlSession.Raw("select LAST_INSERT_ID() as id").Pluck("id", &commentId)
	// 获取创建日期
	commentRet := model.Comment{}
	if err := dao.SqlSession.Model(&model.Comment{}).Where("id=? ", commentId).Find(&commentRet).Error; err != nil {
		// 评论查询失败
		tx.Rollback()
		return 0, "", common.ErrorSelection
	}

	// 提交事务
	tx.Commit()
	//返回评论信息和创建日期
	return commentRet.ID, commentRet.CreatedAt.Format("2006-01-02 15:04:05"), nil

}

// DeleteCommentById 删除评论的持久化操作
func DeleteCommentById(commentId string, userId string, videoId string) (string, string, error) {
	// 开启事务
	tx := dao.SqlSession.Begin()

	// 查询评论信息和创建日期
	commentRet := model.Comment{}
	if err := dao.SqlSession.Model(&model.Comment{}).Where("id=?", commentId).Find(&commentRet).Error; err != nil {
		// 评论查询失败
		tx.Rollback()
		return "", "", common.ErrorSelection
	}

	// 删除评论
	if err := dao.SqlSession.Model(&model.Comment{}).Where("id=?", commentId).Delete(&commentRet).Error; err != nil {
		// 评论删除失败
		tx.Rollback()
		return "", "", common.ErrorDeletion
	}

	// 提交事务
	tx.Commit()
	//返回评论信息和创建日期
	return commentRet.Content, commentRet.CreatedAt.Format("2006-01-02 15:04:05"), nil
}
