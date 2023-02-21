package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name           string `json:"name"`            // 用户名
	Password       string `json:"password"`        // 密码
	FollowCount    uint   `json:"follow_count"`    // 关注人数
	FollowerCount  uint   `json:"follower_count"`  // 粉丝数
	TotalFavorited uint   `json:"total_favorited"` // 总获赞数
	FavoriteCount  uint   `json:"favorite_count"`  // 点赞数
}
