package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"bluebell_backend/pkg/snowflake"
)

// 初始化消费者
func InitConsumer() error {
	// 获取Kafka客户端
	kafkaClient := GetKafkaClient()
	if kafkaClient == nil {
		return fmt.Errorf("kafka client is nil")
	}

	// 创建comment消费者
	go consumeCommentMessages(kafkaClient.client)

	// 创建post消费者
	go consumePostMessages(kafkaClient.client)

	zap.L().Info("kafka consumers started successfully")
	return nil
}

// 消费评论消息
func consumeCommentMessages(client sarama.Client) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		zap.L().Error("new comment consumer failed", zap.Error(err))
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("Comment", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("comment consume partition failed", zap.Error(err))
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// 解析消息
			var comment models.Comment
			if err := json.Unmarshal(msg.Value, &comment); err != nil {
				zap.L().Error("unmarshal comment message failed", zap.Error(err))
				continue
			}

			// 给评论者加积分 (评论+2分)
			if err := addPoints(int64(comment.AuthorID), 2, "评论帖子", "comment", strconv.FormatUint(comment.PostID, 10)); err != nil {
				zap.L().Error("add points for comment author failed", zap.Error(err))
			}

			// 获取帖子作者ID
			// 直接使用uint64类型的postID
			post, err := mysql.GetPostByID(comment.PostID)
			if err != nil {
				fmt.Printf("get post by id failed, postID:%v\n", comment.PostID)
				zap.L().Error("get post by id failed ", zap.Error(err))
				continue
			}

			// 给帖子作者加积分 (被评论+1分)
			if err := addPoints(int64(post.AuthorId), 1, "帖子被评论", "comment", strconv.FormatUint(comment.CommentID, 10)); err != nil {
				zap.L().Error("add points for post author failed", zap.Error(err))
			} else {
				fmt.Printf("add points for post author successfully, postID:%v\n", comment.PostID)
			}

		case err := <-partitionConsumer.Errors():
			zap.L().Error("comment consumer error", zap.Error(err))
		}
	}
}

// 消费帖子消息
func consumePostMessages(client sarama.Client) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		zap.L().Error("new post consumer failed", zap.Error(err))
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("Post-0", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("post consume partition failed", zap.Error(err))
		return
	}
	defer partitionConsumer.Close()
	partitionConsumer10, err := consumer.ConsumePartition("Post-10", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("post consume partition failed", zap.Error(err))
		return
	}
	defer partitionConsumer10.Close()
	partitionConsumer100, err := consumer.ConsumePartition("Post-100", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("post consume partition failed", zap.Error(err))
		return
	}
	defer partitionConsumer100.Close()
	partitionConsumer200, err := consumer.ConsumePartition("Post-200", 0, sarama.OffsetNewest)
	if err != nil {
		zap.L().Error("post consume partition failed", zap.Error(err))
		return
	}
	defer partitionConsumer200.Close()

	for {
		select {
		case msg := <-partitionConsumer200.Messages():
			// 解析消息
			processMessage(msg, "高优先级")
		default:
			select {
			case msg := <-partitionConsumer100.Messages():
				processMessage(msg, "中优先级")
			default:
				select {
				case msg := <-partitionConsumer10.Messages():
					processMessage(msg, "低优先级")
				default:
					select {
					case msg := <-partitionConsumer.Messages():
						processMessage(msg, "超低优先级")
					default:
					}

				}

			}
		}
		select {
		case err := <-partitionConsumer.Errors():
			zap.L().Error("post consumer error", zap.Error(err))
		case err := <-partitionConsumer10.Errors():
			zap.L().Error("post consumer error", zap.Error(err))
		case err := <-partitionConsumer100.Errors():
			zap.L().Error("post consumer error", zap.Error(err))
		case err := <-partitionConsumer200.Errors():
			zap.L().Error("post consumer error", zap.Error(err))
		default:
			// 无错误
		}
	}
}

func processMessage(msg *sarama.ConsumerMessage, priority string) {
	var post models.Post
	if err := json.Unmarshal(msg.Value, &post); err != nil {
		zap.L().Error("unmarshal post message failed",
			zap.String("priority", priority),
			zap.Error(err))
		return
	}

	if err := addPoints(int64(post.AuthorId), 5, "发布帖子", "post", strconv.FormatUint(post.PostID, 10)); err != nil {
		zap.L().Error("add points failed",
			zap.String("priority", priority),
			zap.Uint64("post_id", post.PostID),
			zap.Error(err))
	} else {
		zap.L().Info("points added",
			zap.String("priority", priority),
			zap.Uint64("post_id", post.PostID))
	}
}

// 添加积分
func addPoints(userID int64, points int, description, extType, extID string) error {
	// 1. 查询用户当前积分
	userPoints, err := mysql.GetUserPointsByUserID(userID)
	if err != nil {
		// 用户积分记录不存在，创建新记录
		zap.L().Error("get user points failed", zap.Error(err))
		if errors.Is(err, mysql.ErrUserPointsNotExist) {
			userPoints = &models.UserPoints{
				UserID:      userID,
				Points:      0,
				PointsTotal: 0,
			}
			zap.L().Info("创建用户积分")
		} else {
			return err
		}
	}

	// 2. 更新用户积分
	newPoints := userPoints.Points + points
	newPointsTotal := userPoints.PointsTotal + points

	if err := mysql.UpdateUserPoints(userID, newPoints, newPointsTotal); err != nil {
		return err
	}

	// 3. 记录积分交易
	extJson := fmt.Sprintf(`{"type":"%s","id":"%s"}`, extType, extID)

	// 确定交易类型
	var transactionType int
	switch extType {
	case "post":
		transactionType = 2 // 发布帖子
	case "comment":
		// 根据描述区分是评论还是被评论
		if description == "评论帖子" {
			transactionType = 3 // 评论
		} else if description == "帖子被评论" {
			transactionType = 6 // 被评论 (扩展原有类型)
		}
	}

	// 生成唯一交易ID
	transactionID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("generate transaction ID failed", zap.Error(err))
		return err
	}

	transaction := &models.UserPointsTransaction{
		TransactionID:   int64(transactionID),
		UserID:          userID,
		PointsChange:    points,
		CurrentBalance:  newPoints,
		TransactionType: transactionType,
		Description:     description,
		ExtJson:         extJson,
	}

	if err := mysql.CreateUserPointsTransaction(transaction); err != nil {
		return err
	}

	zap.L().Info("add points successfully",
		zap.Int64("user_id", userID),
		zap.Int("points", points),
		zap.String("description", description))

	return nil
}
