package main

import (
	"checkin_backend/conf"
	"checkin_backend/pb"
	"checkin_backend/pkg/logging"
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net"
	"time"

	"google.golang.org/grpc"
)

// hello server

type server struct {
	pb.UnimplementedCheckinServiceServer
}

const (
	yearSignKeyFormat   = "user:checkins:daily:%d:%d"      // user:checkins:daily:123213131:2025
	monthRetroKeyFormat = "user:checkins:retro:%d:%d:%02d" // user:checkins:retro:123213131:2025:01
)

func (s *server) UserCheckin(ctx context.Context, in *pb.UserCheckinRequest) (*pb.UserCheckinResponse, error) {
	userId := in.GetUserId()
	t := time.Unix(in.GetTimestamp(), 0)
	checkinType := in.GetCheckinType()

	// 检查补签次数
	if checkinType == "补签" {
		retroTimesKey := fmt.Sprintf("user:retro_times:%d", userId)
		remainTimes, err := RedisClient.Get(retroTimesKey).Int64()
		if err != nil || remainTimes <= 0 {
			return &pb.UserCheckinResponse{
				Success:          false,
				Message:          "补签次数不足",
				ContinuousDays:   0,
				TotalDays:        0,
				RemainRetroTimes: 0,
			}, nil
		}
		// 扣除补签次数
		RedisClient.Decr(retroTimesKey)
	}

	year := t.Year()
	key := fmt.Sprintf(yearSignKeyFormat, userId, year)
	offset := t.YearDay() - 1

	// 执行签到
	zap.L().Sugar().Debugf("--> daily setbit key: %s, offset: %d", key, offset)
	ret, err := RedisClient.SetBit(key, int64(offset), 1).Result()
	if err != nil {
		zap.L().Error("daily setbit error", zap.Error(err))
		return &pb.UserCheckinResponse{
			Success:          false,
			Message:          "签到失败",
			ContinuousDays:   0,
			TotalDays:        0,
			RemainRetroTimes: 0,
		}, nil
	}

	// 已签到处理
	if ret == 1 {
		return &pb.UserCheckinResponse{
			Success:          false,
			Message:          "今日已经签到",
			ContinuousDays:   0,
			TotalDays:        0,
			RemainRetroTimes: 0,
		}, nil
	}
	// 计算连续签到天数
	continuousDays := calculateContinuousDays(ctx, userId, t)
	totalDays := calculateTotalDays(ctx, userId)

	return &pb.UserCheckinResponse{
		Success:          true,
		Message:          "签到成功",
		ContinuousDays:   continuousDays,
		TotalDays:        totalDays,
		RemainRetroTimes: getRemainRetroTimes(ctx, userId),
	}, nil
}

func calculateContinuousDays(ctx context.Context, userId uint64, t time.Time) int32 {
	return 1
}
func calculateTotalDays(ctx context.Context, userId uint64) int32 {
	return 1
}

func getRemainRetroTimes(ctx context.Context, userId uint64) int32 {

	return 1
}
func (s *server) GetCheckinStatus(ctx context.Context, in *pb.CheckinStatusRequest) (*pb.CheckinStatusResponse, error) {
	return &pb.CheckinStatusResponse{}, nil
}

var (
	//DB          *gorm.DB
	RedisClient *redis.Client
)

// MustInitMySQL 初始化 MySQL 连接
//func MustInitMySQL(cfg *viper.Viper) {
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//		cfg.GetString("mysql.user"),
//		cfg.GetString("mysql.password"),
//		cfg.GetString("mysql.host"),
//		cfg.GetString("mysql.port"),
//		cfg.GetString("mysql.dbname"),
//	)
//	db, err := gorm.Open(mysql.Open(dsn))
//	if err != nil {
//		panic(fmt.Errorf("connect db fail: %w", err))
//	}
//
//	sqlDB, err := db.DB()
//	if err != nil {
//		panic(fmt.Errorf("connect db fail: %w", err))
//	}
//	// 设置连接池参数
//	sqlDB.SetMaxIdleConns(cfg.GetInt("mysql.max_idle_conns"))
//	sqlDB.SetMaxOpenConns(cfg.GetInt("mysql.max_open_conns"))
//	sqlDB.SetConnMaxLifetime(cfg.GetDuration("mysql.max_lifetime"))
//
//	query.SetDefault(db) // 指定 query 包使用的默认数据库连接
//}

// MustInitRedis 初始化 Redis 连接
func MustInitRedis(conf *viper.Viper) {
	addr := fmt.Sprintf("%s:%d", conf.GetString("redis.host"), conf.GetInt("redis.port"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.GetString("redis.password"),
		DB:       conf.GetInt("redis.db"),
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		panic(fmt.Errorf("init redis failed, err:%w", err))
	}
	RedisClient = rdb
}

var confPath = flag.String("conf", "./config/config.yaml", "配置文件路径")

func main() {
	flag.Parse()
	cfg := conf.Load(*confPath)
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer logger.Sync()
	//MustInitMySQL(cfg)
	MustInitRedis(cfg)
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	srv := &server{}
	s := grpc.NewServer()                   // 创建gRPC服务器
	pb.RegisterCheckinServiceServer(s, srv) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
