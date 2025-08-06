package kafka

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type KafkaClient struct {
	config   *sarama.Config
	client   sarama.Client
	producer sarama.AsyncProducer
}

func NewKafkaClient() (*KafkaClient, error) {
	config := sarama.NewConfig()
	// 生产者配置
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Return.Errors = true // 启用错误通道
	// 消费者配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest         // 从最新消息开始消费
	config.Consumer.Offsets.AutoCommit.Enable = true              // 启用自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动提交间隔
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		zap.L().Error("new kafka client failed", zap.Error(err))
		return nil, err
	}
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		zap.L().Error("new kafka producer failed", zap.Error(err))
		return nil, err
	}
	return &KafkaClient{
		config:   config,
		client:   client,
		producer: producer,
	}, nil
}

func (k *KafkaClient) SendMessage(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
		Key:   sarama.ByteEncoder(key),
	}
	k.producer.Input() <- msg

	// 等待生产者响应
	select {
	case success := <-k.producer.Successes():
		zap.L().Info("Message sent successfully",
			zap.String("topic", success.Topic),
			zap.Int32("partition", success.Partition),
			zap.Int64("offset", success.Offset))
		return nil
	case err := <-k.producer.Errors():
		zap.L().Error("Failed to send message",
			zap.String("topic", topic),
			zap.Error(err))
		return err.Err
	}
}
func (k *KafkaClient) Close() error {
	if k.producer != nil {
		if err := k.producer.Close(); err != nil {
			zap.L().Error("close kafka producer failed", zap.Error(err))
			return err
		}
	}
	if k.client != nil {
		if err := k.client.Close(); err != nil {
			zap.L().Error("close kafka client failed", zap.Error(err))
			return err
		}
	}
	zap.L().Info("kafka client closed successfully")
	return nil
}

var client *KafkaClient

func InitKafka() error {
	var err error
	client, err = NewKafkaClient()
	if err != nil {
		return err
	}

	// 初始化消费者
	if err = InitConsumer(); err != nil {
		zap.L().Error("init kafka consumer failed", zap.Error(err))
		return err
	}

	return nil
}
func GetKafkaClient() *KafkaClient {
	return client
}

// CloseKafka 关闭Kafka客户端
func CloseKafka() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
