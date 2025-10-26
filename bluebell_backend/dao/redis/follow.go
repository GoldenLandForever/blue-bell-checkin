package redis

import (
	"strconv"
)

// Follow 关注某人
func Follow(userId, followedId uint64) error {
	// 开启事务
	pipeline := Client.TxPipeline()

	// 将被关注者添加到用户的关注列表
	followKey := KeyFollowSetPrefix + strconv.FormatUint(userId, 10)
	pipeline.SAdd(followKey, followedId)

	// 将关注者添加到被关注者的粉丝列表
	followerKey := KeyFollowerSetPrefix + strconv.FormatUint(followedId, 10)
	pipeline.SAdd(followerKey, userId)

	_, err := pipeline.Exec()
	return err
}

// UnFollow 取消关注
func UnFollow(userId, followedId uint64) error {
	// 开启事务
	pipeline := Client.TxPipeline()

	// 从用户的关注列表中移除被关注者
	followKey := KeyFollowSetPrefix + strconv.FormatUint(userId, 10)
	pipeline.SRem(followKey, followedId)

	// 从被关注者的粉丝列表中移除关注者
	followerKey := KeyFollowerSetPrefix + strconv.FormatUint(followedId, 10)
	pipeline.SRem(followerKey, userId)

	_, err := pipeline.Exec()
	return err
}

// IsFollowing 检查是否已关注
func IsFollowing(userId, followedId uint64) (bool, error) {
	followKey := KeyFollowSetPrefix + strconv.FormatUint(userId, 10)
	return Client.SIsMember(followKey, followedId).Result()
}

// GetFollowings 获取关注列表
func GetFollowings(userId uint64) ([]string, error) {
	followKey := KeyFollowSetPrefix + strconv.FormatUint(userId, 10)
	return Client.SMembers(followKey).Result()
}

// GetFollowers 获取粉丝列表
func GetFollowers(userId uint64) ([]string, error) {
	followerKey := KeyFollowerSetPrefix + strconv.FormatUint(userId, 10)
	return Client.SMembers(followerKey).Result()
}
