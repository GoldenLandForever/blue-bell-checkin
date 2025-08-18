package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"checkin_backend/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server 实现了CheckinServiceServer接口
type server struct {
	pb.UnimplementedCheckinServiceServer
}

// UserCheckin 实现用户签到功能
func (s *server) UserCheckin(ctx context.Context, in *pb.UserCheckinRequest) (*pb.UserCheckinResponse, error) {
	log.Printf("收到签到请求: user_id=%d, timestamp=%d, checkin_type=%s", 
		in.UserId, in.Timestamp, in.CheckinType)
	
	// TODO: 这里可以添加实际的签到逻辑
	// 示例响应
	return &pb.UserCheckinResponse{
		Success:         true,
		Message:         "签到成功",
		ContinuousDays:  1,
		TotalDays:       1,
		RemainRetroTimes: 3,
	}, nil
}

// GetCheckinStatus 实现获取签到状态功能
func (s *server) GetCheckinStatus(ctx context.Context, in *pb.CheckinStatusRequest) (*pb.CheckinStatusResponse, error) {
	log.Printf("收到状态查询请求: user_id=%d, timestamp=%d", in.UserId, in.Timestamp)
	
	// TODO: 这里可以添加实际的状态查询逻辑
	// 示例响应
	return &pb.CheckinStatusResponse{
		CheckedinDays:      []int32{1, 5, 15, 20, 25}, // 当月已签到日期
		RetroCheckedinDays: []int32{3, 12, 18},        // 当月补签日期
	}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	// 创建gRPC服务器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor), // 添加日志中间件
	)
	
	// 注册服务
	pb.RegisterCheckinServiceServer(s, &server{})
	
	// 注册反射服务，方便调试
	reflection.Register(s)
	
	log.Printf("gRPC server starting on :8972...")
	
	// 启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// loggingInterceptor 日志中间件
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s, Duration: %v, Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}