package controller

import (
	"bluebell_backend/logic"
	"bluebell_backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FollowHandler 关注用户
func FollowHandler(c *gin.Context) {
	// 获取当前用户ID
	uidStr, exists := c.Get("userID")
	if !exists {
		ResponseError(c, CodeNotLogin)
		return
	}
	userId := uidStr.(uint64)

	// 解析请求参数
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		zap.L().Error("FollowHandler ShouldBindJSON failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 不能关注自己
	if userId == follow.FollowedID {
		zap.L().Error("cannot follow yourself")
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 检查是否已经关注
	isFollowing, err := logic.GetFollowStatus(userId, follow.FollowedID)
	if err != nil {
		zap.L().Error("logic.GetFollowStatus failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	if isFollowing {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 执行关注操作
	if err := logic.Follow(userId, &follow); err != nil {
		zap.L().Error("logic.Follow failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// UnFollowHandler 取消关注
func UnFollowHandler(c *gin.Context) {
	// 获取当前用户ID
	uidStr, exists := c.Get("userID")
	if !exists {
		ResponseError(c, CodeNotLogin)
		return
	}
	userId := uidStr.(uint64)

	// 解析请求参数
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		zap.L().Error("UnFollowHandler ShouldBindJSON failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 不能取消关注自己
	if userId == follow.FollowedID {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 检查是否已经关注
	isFollowing, err := logic.GetFollowStatus(userId, follow.FollowedID)
	if err != nil {
		zap.L().Error("logic.GetFollowStatus failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	if !isFollowing {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 执行取消关注操作
	if err := logic.UnFollow(userId, &follow); err != nil {
		zap.L().Error("logic.UnFollow failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// GetFollowListHandler 获取关注列表
func GetFollowListHandler(c *gin.Context) {
	// 获取用户ID（从URL参数或当前登录用户）
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取关注列表
	followings, err := logic.GetFollowings(userId)
	if err != nil {
		zap.L().Error("logic.GetFollowings failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, followings)
}

// GetFollowerListHandler 获取粉丝列表
func GetFollowerListHandler(c *gin.Context) {
	// 获取用户ID（从URL参数或当前登录用户）
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取粉丝列表
	followers, err := logic.GetFollowers(userId)
	if err != nil {
		zap.L().Error("logic.GetFollowers failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, followers)
}
