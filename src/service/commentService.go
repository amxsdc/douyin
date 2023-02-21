package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"strconv"
)

// AddComment 添加评论的持久化操作
// TODO 要使用声明式事务
func AddComment(commentText string, userId string, videoId string) (uint, string, error) {
	// 添加评论
	//将string类型数据转化为uint类型数据
	vId, _ := strconv.Atoi(videoId)
	uId, _ := strconv.Atoi(userId)
	comment := model.Comment{
		VideoId: uint(vId),
		UserId:  uint(uId),
		Content: commentText,
	}
	dao.SqlSession.Model(&model.Comment{}).Create(&comment)
	// 获取评论id
	var commentId []uint
	dao.SqlSession.Raw("select LAST_INSERT_ID() as id").Pluck("id", &commentId)
	// 获取创建日期
	commentRet := model.Comment{}
	dao.SqlSession.Model(&model.Comment{}).Where("id=? ", commentId).Find(&commentRet)

	//返回评论信息和创建日期
	return commentRet.ID, commentRet.CreatedAt.Format("2006-01-02 15:04:05"), nil

}

// DeleteCommentById 删除评论的持久化操作
// TODO 要使用声明式事务
func DeleteCommentById(commentId string, userId string, videoId string) (string, string, error) {

	// 查询评论信息和创建日期
	commentRet := model.Comment{}
	dao.SqlSession.Model(&model.Comment{}).Where("id=?", commentId).Find(&commentRet)

	// 删除评论
	dao.SqlSession.Model(&model.Comment{}).Where("id=?", commentId).Delete(&commentRet)

	//返回评论信息和创建日期
	return commentRet.Content, commentRet.CreatedAt.Format("2006-01-02 15:04:05"), nil
}
