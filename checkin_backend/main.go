package main

import (
	"checkin_backend/pb"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// hello server

type server struct {
	pb.UnimplementedCheckinServiceServer
}

func (s *server) UserCheckin(ctx context.Context, in *pb.UserCheckinRequest) (*pb.UserCheckinResponse, error) {
	return &pb.UserCheckinResponse{
		Success:          true,
		Message:          "签到成功",
		ContinuousDays:   1,
		TotalDays:        1,
		RemainRetroTimes: 3,
	}, nil
}

func (s *server) GetCheckinStatus(ctx context.Context, in *pb.CheckinStatusRequest) (*pb.CheckinStatusResponse, error) {
	return &pb.CheckinStatusResponse{}, nil
}

func main() {
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
