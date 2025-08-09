package tests

import (
	"bluebell_backend/dao/redis"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPostListPerformance 测试帖子列表性能
func TestPostListPerformance(t *testing.T) {
	const (
		concurrentUsers = 1000
		duration        = 30 * time.Second
	)

	var wg sync.WaitGroup
	start := time.Now()
	var successCount int64
	var errorCount int64

	// 启动并发测试
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			endTime := time.Now().Add(duration)
			for time.Now().Before(endTime) {
				// 模拟获取帖子列表
				startReq := time.Now()
				
				// 从Redis获取最新帖子
				postIDs, err := redis.client.ZRevRange("post:new", 0, 9).Result()
				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					continue
				}

				// 获取帖子详情
				for _, postID := range postIDs {
					_, err := redis.Client.HGetAll(fmt.Sprintf("post:detail:%s", postID)).Result()
					if err != nil {
						atomic.AddInt64(&errorCount, 1)
						break
					}
				}

				elapsed := time.Since(startReq)
				if elapsed < 50*time.Millisecond {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)

	// 计算性能指标
	qps := float64(successCount) / totalTime.Seconds()
	errorRate := float64(errorCount) / float64(successCount+errorCount) * 100

	t.Logf("并发用户数: %d", concurrentUsers)
	t.Logf("总请求数: %d", successCount+errorCount)
	t.Logf("成功请求: %d", successCount)
	t.Logf("失败请求: %d", errorCount)
	t.Logf("QPS: %.2f", qps)
	t.Logf("错误率: %.2f%%", errorRate)

	// 断言性能指标
	assert.Greater(t, qps, 1000.0, "QPS应该大于1000")
	assert.Less(t, errorRate, 1.0, "错误率应该小于1%")
}

// TestVotePerformance 测试投票性能
func TestVotePerformance(t *testing.T) {
	const (
		concurrentUsers = 500
		duration        = 30 * time.Second
	)

	var wg sync.WaitGroup
	start := time.Now()
	var voteCount int64

	// 模拟投票操作
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			endTime := time.Now().Add(duration)
			for time.Now().Before(endTime) {
				// 模拟用户投票
				userID := uint64(time.Now().UnixNano())
				postID := uint64(time.Now().Unix() % 1000) // 随机帖子
				
				// 检查投票记录
				voteKey := fmt.Sprintf("vote:%d:%d", userID, postID)
				exists, _ := redis.Client.Exists(voteKey).Result()
				
				if exists == 0 {
					// 记录投票
					pipe := redis.Client.Pipeline()
					pipe.Set(voteKey, 1, 7*24*time.Hour)
					pipe.ZIncrBy(fmt.Sprintf("post:votes:%d", postID), 1, fmt.Sprintf("user:%d", userID))
					_, _ = pipe.Exec()
					
					atomic.AddInt64(&voteCount, 1)
				}
				
				time.Sleep(10 * time.Millisecond) // 控制频率
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)
	voteQPS := float64(voteCount) / totalTime.Seconds()

	t.Logf("投票操作总数: %d", voteCount)
	t.Logf("投票QPS: %.2f", voteQPS)

	assert.Greater(t, voteQPS, 100.0, "投票QPS应该大于100")
}

// BenchmarkRedisCache 测试Redis缓存性能
func BenchmarkRedisCache(b *testing.B) {
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// 模拟缓存访问模式
		postID := uint64(i % 1000)
		
		// 1. 检查缓存
		cacheKey := fmt.Sprintf("post:detail:%d", postID)
		_, err := redis.Client.Exists(cacheKey).Result()
		if err != nil {
			b.Fatal(err)
		}
		
		// 2. 获取数据
		_, err = redis.Client.HGetAll(cacheKey).Result()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestMemoryUsage 测试内存使用情况
func TestMemoryUsage(t *testing.T) {
	// 预热数据
	warmup := redis.NewRedisWarmUp()
	err := warmup.WarmUpAll()
	assert.NoError(t, err)
	
	// 等待数据加载
	time.Sleep(2 * time.Second)
	
	// 获取Redis内存信息
	info, err := redis.Client.Info("memory").Result()
	assert.NoError(t, err)
	
	t.Logf("Redis内存信息:\n%s", info)
	
	// 检查关键指标
	usedMemory, _ := redis.Client.Info("memory").Result()
	t.Logf("已使用内存: %s", usedMemory)
}