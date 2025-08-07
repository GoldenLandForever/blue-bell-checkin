package redis

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"bluebell_backend/settings"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

type SliceCmd = redis.SliceCmd
type StringStringMapCmd = redis.StringStringMapCmd

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func SyncExpiringVotes(expire int) (err error) {
	// 从Redis中获取所有过期的投票键
	iter := client.Scan(0, "*voted:*", 1000).Iterator()
	if err != nil {
		return err
	}

	// 遍历所有键，检查是否过期
	for iter.Next() {
		key := iter.Val()
		// 检查键是否过期
		ttl, err := client.TTL(key).Result()
		if err != nil {
			return err
		}

		// 如果过期时间小于当前时间，说明键已过期
		if ttl <= 0 {
			// 从Redis中删除过期键
			if err := client.Del(key).Err(); err != nil {
				continue
			}
		}
		// 将还剩下7天过期的key持久化到mysql中
		if ttl < time.Duration(expire)*time.Second {
			// 从key中提取投票信息
			go persistVotes(key)
		}
	}

	return nil
}

func persistVotes(key string) {
	// 分页获取ZSet数据
	// MySQL事务使用ON DUPLICATE KEY UPDATE
	// 成功持久化后删除Redis键
	postID := strings.Split(key, "voted:")[1]
	// 从Redis中获取投票信息
	members, err := client.ZRangeWithScores(key, 0, -1).Result()
	if err != nil {
		zap.L().Error("获取ZSet成员失败", zap.Error(err))
		return
	}
	tx, _ := mysql.Db.Begin()
	for _, member := range members {
		if userID, ok := member.Member.(string); ok {
			userIDInt, _ := strconv.ParseUint(userID, 10, 64)
			postIDInt, _ := strconv.ParseUint(postID, 10, 64)
			record := models.Vote{
				UserID:     userIDInt,
				PostID:     postIDInt,
				VoteType:   int8(member.Score),
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			sqlstr := `insert into vote_records (user_id, post_id, vote_type,created_at,updated_at) values (?, ?, ?,?,?)`
			if _, err := tx.Exec(sqlstr, record.UserID, record.PostID, record.VoteType, record.CreateTime, record.UpdateTime); err != nil {
				zap.L().Error("执行失败回滚", zap.Error(err))
				tx.Rollback()
				break
			}
		}

	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("事务提交失败", zap.Error(err))
	}
}

func Close() {
	_ = client.Close()
}
