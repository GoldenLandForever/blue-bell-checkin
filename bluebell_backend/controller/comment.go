package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"bluebell_backend/pkg/kafka"
	"bluebell_backend/pkg/snowflake"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 评论

// CommentHandler 创建评论
func CommentHandler(c *gin.Context) {
	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		fmt.Println(err)
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 生成评论ID
	commentID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 获取作者ID，当前请求的UserID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	comment.CommentID = commentID
	comment.AuthorID = userID
	// post_id已通过json标签自动从字符串转换为uint64
	if comment.PostID == 0 {
		zap.L().Error("post_id绑定失败")
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 创建评论
	if err := mysql.CreateComment(&comment); err != nil {
		zap.L().Error("mysql.CreateComment(&comment) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
	// 创建kafka积分生产者

	// 发送评论消息
	kafkaClient := kafka.GetKafkaClient()
	// 将积分变动对象序列化为JSON
	scoreData, err := json.Marshal(comment)
	if err != nil {
		zap.L().Error("序列化积分消息失败", zap.Error(err))
		return
	}
	// 发送积分消息到Kafka
	if err := kafkaClient.SendMessage("Comment", []byte(fmt.Sprintf("%d", userID)), scoreData); err != nil {
		zap.L().Error("发送积分消息失败", zap.Error(err))
	}

}

// CommentListHandler 评论列表
func CommentListHandler(c *gin.Context) {
	//输入postID
	//mysql在comment表里面查询postID，暂且按照时间顺序排序
	var commentlist models.CommentList
	if err := c.ShouldBindQuery(&commentlist); err != nil {
		fmt.Println(err)
		ResponseError(c, CodeInvalidParams)
		return
	}
	posts, err := mysql.GetCommentListByIDs(commentlist)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("发送评论信息成功")
	ResponseSuccess(c, posts)
}
