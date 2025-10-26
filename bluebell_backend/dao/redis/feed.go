package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// PushFeed 推送帖子到用户的feed流中
func PushFeed(userID uint64, postID uint64) error {
	// 使用 ZADD 命令将帖子添加到用户的feed流中
	feedKey := KeyFeedPrefix + strconv.FormatUint(userID, 10)
	return Client.ZAdd(feedKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.FormatUint(postID, 10),
	}).Err()
}

// PushFeedToFollowers 推送帖子到所有在线粉丝的feed流中
func PushFeedToFollowers(authorID uint64, postID uint64) error {
	// 1. 获取作者的粉丝列表
	followerKey := KeyFollowerSetPrefix + strconv.FormatUint(authorID, 10)
	followers, err := Client.SMembers(followerKey).Result()
	if err != nil {
		return err
	}

	// 2. 遍历粉丝列表
	for _, followerStr := range followers {
		followerID, err := strconv.ParseUint(followerStr, 10, 64)
		if err != nil {
			continue
		}

		// 3. 检查粉丝是否在线
		isOnline, err := GetUserOnlineStatus(followerID)
		if err != nil || !isOnline {
			continue
		}

		// 4. 推送到在线粉丝的feed流中
		if err := PushFeed(followerID, postID); err != nil {
			zap.L().Error("PushFeed to follower failed", zap.Uint64("follower_id", followerID), zap.Error(err))
			continue // 继续处理下一个粉丝，不中断整个过程
		}
	}
	return nil
}

// GetUserFeed 获取用户的feed流
func GetUserFeed(userID uint64, offset, limit int64) ([]string, error) {
	// 使用 ZREVRANGE 命令获取用户的feed流，按时间倒序排列
	feedKey := KeyFeedPrefix + strconv.FormatUint(userID, 10)
	return Client.ZRevRange(feedKey, offset, offset+limit-1).Result()
}

// RemoveExpiredFeed 删除过期的feed条目（可选：保留最近7天或最新1000条）
func RemoveExpiredFeed(userID uint64) error {
	feedKey := KeyFeedPrefix + strconv.FormatUint(userID, 10)

	// 方案1：删除7天前的数据
	old := time.Now().AddDate(0, 0, -7).Unix()
	if err := Client.ZRemRangeByScore(feedKey, "0", strconv.FormatInt(old, 10)).Err(); err != nil {
		return err
	}

	// 方案2：只保留最新的1000条
	count, err := Client.ZCard(feedKey).Result()
	if err != nil {
		return err
	}
	if count > 1000 {
		if err := Client.ZRemRangeByRank(feedKey, 0, count-1001).Err(); err != nil {
			return err
		}
	}

	return nil
}
