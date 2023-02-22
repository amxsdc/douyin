package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"strconv"
)

// GetVideoById 根据视频id获取视频信息
func GetVideoById(videoId string) (model.Video, error) {
	var video = model.Video{}
	if err := dao.SqlSession.Model(&model.Video{}).Where("id=?", videoId).Find(&video).Error; err != nil {
		return model.Video{}, common.ErrorSelection
	}
	return video, nil
}

// UpdateFavoriteCount 更新视频点赞数
func UpdateFavoriteCount(video *model.Video, user *model.User, author *model.User, change int) error {
	// 上锁
	m.Lock()

	// 设置解锁
	defer m.Unlock()

	// 判断是否为负
	if video.FavoriteCount == 0 && change < 0 {
		return common.ErrorFavoriteUpdate
	}

	if change < 0 {
		// 视频点赞数
		video.FavoriteCount = video.FavoriteCount - uint(change*-1)
		// 用户点赞数
		user.FavoriteCount = user.FavoriteCount - uint(change*-1)
		// 作者被点赞总数
		author.TotalFavorited = author.TotalFavorited - uint(change*-1)
	} else {
		// 视频点赞数
		video.FavoriteCount = video.FavoriteCount + uint(change)
		// 用户点赞数
		user.FavoriteCount = user.FavoriteCount + uint(change)
		// 作者被点赞总数
		author.TotalFavorited = author.TotalFavorited + uint(change)
	}

	return nil
}

// FavoriteAction 点赞/取消点赞行为
func FavoriteAction(videoId string, userId uint, count int) error {
	// 开启事务
	tx := dao.SqlSession.Begin()

	// 根据videoId获取视频信息
	video, err := GetVideoById(videoId)
	if err != nil {
		tx.Rollback()
		return common.ErrorSelection
	}
	// 获取视频作者信息用户更新
	author, err := GetUser(video.AuthorId)
	if err != nil {
		tx.Rollback()
		return common.ErrorUpdate
	}
	// 获取用户以便进行更新
	user, err := GetUser(userId)
	if err != nil {
		tx.Rollback()
		return common.ErrorUpdate
	}

	// 判断是否重复点赞
	if IsFavoriteRepeated(userId, videoId) && count > 0 {
		return common.ErrorFavoriteRepeat
	}

	// 判断是否为自己给自己点赞
	if user.ID == author.ID {
		err = UpdateFavoriteCount(&video, &user, &user, count)
		dao.SqlSession.Save(video)
		dao.SqlSession.Save(user)
	} else {
		// 更新点赞信息
		err = UpdateFavoriteCount(&video, &user, &author, count)

		// 保存相关信息
		dao.SqlSession.Save(video)
		dao.SqlSession.Save(author)
		dao.SqlSession.Save(user)
	}

	if err != nil {
		tx.Rollback()
		return common.ErrorUpdate
	}

	vId, _ := strconv.Atoi(videoId)

	// 增添或删除记录
	if count > 0 {
		err = InsertFavorite(uint(vId), userId)
	} else if count < 0 {
		err = DeleteFavorite(uint(vId), userId)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}

// IsFavoriteRepeated 软删除时使用 判断表中是否有该id下的信息
func IsFavoriteRepeated(userId uint, videoId string) bool {
	var favorite model.Favorite
	if err := dao.SqlSession.Model(model.Favorite{}).Where("user_id=? and video_id=?", userId, videoId).Find(&favorite).Error; err != nil {
		return false
	}
	return true
}

// InsertFavorite 插入点赞信息
func InsertFavorite(videoId uint, userId uint) error {
	var favorite = model.Favorite{
		VideoId: videoId,
		UserId:  userId,
		State:   0,
	}

	if err := dao.SqlSession.Model(&model.Favorite{}).Create(&favorite).Error; err != nil {
		return common.ErrorFavoriteRepeat
	}
	return nil
}

// DeleteFavorite 删除点赞信息
func DeleteFavorite(videoId uint, userId uint) error {
	var favorite model.Favorite
	if err := dao.SqlSession.Where("video_id = ? and user_id = ?", videoId, userId).Delete(&favorite).Error; err != nil {
		return common.ErrorDeletion
	}
	return nil
}

// GetVideoByUserId 根据userId获取视频信息
func GetVideoByUserId(userId string) ([]model.Video, error) {
	var videoIds []model.Favorite
	// 查询userId对应的video信息
	if err := dao.SqlSession.Model(&model.Favorite{}).Where("user_id=?", userId).Find(&videoIds).Error; err != nil {
		return []model.Video{}, nil
	}

	var videos []model.Video
	// 将视频数组插入到响应体中
	for _, videoId := range videoIds {
		var video model.Video
		dao.SqlSession.Model(&model.Video{}).Where("id=?", videoId.VideoId).Find(&video)
		videos = append(videos, video)
	}

	return videos, nil
}
