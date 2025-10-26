package controller

import (
	"bluebell_backend/dao/redis"
	"bluebell_backend/logic"
	"bluebell_backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserFeedHandler 获取用户的feed流
func GetUserFeedHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetUserFeedHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	//获取离线时间的feed流
	err = redis.PullOfflineFeed(userID)
	if err != nil {
		zap.L().Error("redis.PullOfflineFeed failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 获取feed流中的帖子ID列表
	offset := (p.Page - 1) * p.Size
	postIDs, err := redis.GetUserFeed(userID, int64(offset), int64(p.Size))
	if err != nil {
		zap.L().Error("redis.GetUserFeed failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 获取帖子详情
	posts := make([]*models.ApiPostDetail, 0, len(postIDs))
	for _, pidStr := range postIDs {
		pid, err := strconv.ParseUint(pidStr, 10, 64)
		if err != nil {
			continue
		}
		post, err := logic.GetPostById(pid)
		if err != nil {
			zap.L().Error("logic.GetPostById failed",
				zap.Uint64("post_id", pid),
				zap.Error(err))
			continue
		}
		posts = append(posts, post)
	}

	// 4. 返回响应
	ResponseSuccess(c, posts)
}
