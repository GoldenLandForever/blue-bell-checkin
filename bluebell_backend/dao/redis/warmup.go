package redis

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// RedisWarmUp Redis预热器
type RedisWarmUp struct {
}

// NewRedisWarmUp 创建新的预热器
func NewRedisWarmUp() *RedisWarmUp {
	// 使用默认配置
	return &RedisWarmUp{}
}

// WarmUpAll 执行完整的Redis预热流程
func (w *RedisWarmUp) WarmUpAll() error {
	zap.L().Info("开始Redis预热...")

	start := time.Now()

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// 1. 预热最新帖子数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := w.warmUpNewPostData(); err != nil {
			errChan <- fmt.Errorf("投票数据预热失败: %w", err)
		}
	}()

	// 等待所有预热任务完成
	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		zap.L().Warn("Redis预热部分失败", zap.Any("errors", errors))
		// 不返回错误，让应用继续启动
	}

	elapsed := time.Since(start)
	zap.L().Info("Redis预热完成", zap.Duration("耗时", elapsed))

	return nil
}

func (w *RedisWarmUp) warmUpNewPostData() error {
	sqlstr := `
		SELECT post_id,title, content, author_id, community_id, create_time 
		FROM post 
		ORDER BY create_time DESC 
		LIMIT 1000
	`
	rows, err := mysql.Db.Query(sqlstr)
	if err != nil {
		return fmt.Errorf("查询最新帖子数据失败: %w", err)
	}
	defer rows.Close()
	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.AuthorId, &post.CommunityID, &post.CreateTime); err != nil {
			zap.L().Warn("扫描帖子数据失败", zap.Error(err))
			continue
		}
		posts = append(posts, &post)
	}
	for _, post := range posts {
		// 1. 添加到新帖子有序集合（保持原有功能）
		err := Client.ZAdd(KeyPostNew, redis.Z{
			Score:  float64(post.CreateTime.Unix()),
			Member: post.PostID,
		}).Err()
		if err != nil {
			zap.L().Warn("添加帖子到新帖子集合失败", zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			continue
		}
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "community_id", post.CommunityID)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "author_id", post.AuthorId)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "create_time", post.CreateTime)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "title", post.Title)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "content", post.Content)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "community_name", community.CommunityName)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "author_name", user.UserName)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "introduction", community.Introduction)
		Client.HSet(KeyPostInfoHashPrefix+fmt.Sprint(post.PostID), "create_time", community.CreateTime)
	}
	return nil
}
