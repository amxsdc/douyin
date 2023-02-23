package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserAttr struct {
	Id              uint   `json:"id"`
	Name            string `json:"name"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  uint   `json:"total_favorited"`
	WorkCount       uint   `json:"work_count"`
	FavoriteCount   uint   `json:"favorite_count"`
}

type Follow struct {
	gorm.Model
	HostId  uint
	GuestId uint
}

func ObtainFollowingList(userId uint) ([]UserAttr, error) {

	var FollowIdList []Follow
	dao.SqlSession.AutoMigrate(&model.Following{})
	if err := dao.SqlSession.Table("followings").Where("host_id = ?", userId).Find(&FollowIdList).Error; err != nil {
		return nil, err
	}

	return tidyData(FollowIdList)

}

func ObtainFollowerList(userId uint) ([]UserAttr, error) {

	var FollowerIdList []Follow
	dao.SqlSession.AutoMigrate(&model.Followers{})
	if err := dao.SqlSession.Table("followers").Where("host_id = ?", userId).Find(&FollowerIdList).Error; err != nil {
		return nil, err
	}
	return tidyData(FollowerIdList)

}

func tidyData(FollowIdList []Follow) ([]UserAttr, error) {
	var UserAttrList []UserAttr

	for _, v := range FollowIdList {
		var user model.User
		if err := dao.SqlSession.Model(&model.User{}).Where("id = ?", v.GuestId).First(&user).Error; err != nil {
			return nil, fmt.Errorf("query data failture!:%v", err)
		}
		UserAttrs := UserAttr{
			Id:              v.GuestId,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        true,
			Avatar:          "",
			BackgroundImage: "",
			Signature:       "",
			TotalFavorited:  user.TotalFavorited,
			FavoriteCount:   user.FavoriteCount,
		}
		UserAttrList = append(UserAttrList, UserAttrs)
	}
	return UserAttrList, nil
}
