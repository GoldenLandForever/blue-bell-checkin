package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 设置用户的在线状态过期时间为2分钟
const UserOnlineExpiration = 2 * time.Minute

// UpdateUserOnline 更新用户的在线状态
func UpdateUserOnline(userID uint64) error {
	// 使用 ZADD 命令将用户ID和当前时间戳添加到有序集合中
	return Client.ZAdd(KeyOnlineZSet, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.FormatUint(userID, 10),
	}).Err()
}

// GetUserOnlineStatus 获取用户的在线状态
func GetUserOnlineStatus(userID uint64) (bool, error) {
	userIDStr := strconv.FormatUint(userID, 10)
	// 获取用户的最后心跳时间
	score, err := Client.ZScore(KeyOnlineZSet, userIDStr).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// 检查最后心跳时间是否在有效期内
	lastHeartbeat := time.Unix(int64(score), 0)
	if time.Since(lastHeartbeat) > UserOnlineExpiration {
		// 心跳无效，删除用户的在线状态
		err = Client.ZRem(KeyOnlineZSet, userIDStr).Err()
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// GetOnlineUsers 获取所有在线用户
func GetOnlineUsers() ([]string, error) {
	// 获取当前时间戳
	now := time.Now().Unix()
	// 获取在线用户的时间范围：当前时间往前推2分钟
	min := now - int64(UserOnlineExpiration.Seconds())

	// 使用 ZRANGEBYSCORE 命令获取在线用户
	return Client.ZRangeByScore(KeyOnlineZSet, redis.ZRangeBy{
		Min: strconv.FormatInt(min, 10),
		Max: strconv.FormatInt(now, 10),
	}).Result()
}

// CleanupOfflineUsers 清理离线用户
func CleanupOfflineUsers() error {
	// 获取当前时间戳
	now := time.Now().Unix()
	// 清理超过2分钟没有心跳的用户
	min := now - int64(UserOnlineExpiration.Seconds())

	// 使用 ZREMRANGEBYSCORE 命令删除超时的用户
	return Client.ZRemRangeByScore(KeyOnlineZSet, "-inf", strconv.FormatInt(min, 10)).Err()
}
