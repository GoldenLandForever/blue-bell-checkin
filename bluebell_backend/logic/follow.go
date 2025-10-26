package logic

import (
	"bluebell_backend/dao/redis"
	"bluebell_backend/models"
)

// Follow 关注
func Follow(userId uint64, follow *models.Follow) error {
	return redis.Follow(userId, follow.FollowedID)
}

// UnFollow 取消关注
func UnFollow(userId uint64, follow *models.Follow) error {
	return redis.UnFollow(userId, follow.FollowedID)
}

// GetFollowStatus 获取关注状态
func GetFollowStatus(userId uint64, targetId uint64) (bool, error) {
	return redis.IsFollowing(userId, targetId)
}

// GetFollowings 获取关注列表
func GetFollowings(userId uint64) ([]string, error) {
	return redis.GetFollowings(userId)
}

// GetFollowers 获取粉丝列表
func GetFollowers(userId uint64) ([]string, error) {
	return redis.GetFollowers(userId)
}
