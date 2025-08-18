# 🎯 签到模块使用指南

## 📋 Proto定义说明

我们设计了一个完整的签到模块，包含以下功能：

### ✅ 核心功能
1. **用户签到** - 记录用户每日签到
2. **签到状态查询** - 查看今日是否已签到
3. **签到记录查询** - 分页获取历史签到记录
4. **排行榜** - 查看连续签到/总签到/月度排行榜

### 📊 数据结构优化

相比原始定义，我们增加了：
- **防刷机制** - device_id 和 ip_address
- **奖励系统** - reward_points 和连续签到奖励
- **时间记录** - 精确到时间戳
- **签到类型** - 支持补签、活动签到等
- **分页查询** - 支持历史记录分页

## 🚀 使用示例

### 1. 生成Go代码
```bash
# 在项目根目录执行
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative bluebell_backend/pb/checkin.proto
```

### 2. 服务端实现示例
```go
package main

import (
    "context"
    "log"
    "net"
    "time"

    pb "bluebell_backend/pb"
    "google.golang.org/grpc"
)

type checkinServer struct {
    pb.UnimplementedCheckinServiceServer
}

func (s *checkinServer) UserCheckin(ctx context.Context, req *pb.UserCheckinRequest) (*pb.UserCheckinResponse, error) {
    // TODO: 实现签到逻辑
    return &pb.UserCheckinResponse{
        Success:         true,
        Message:         "签到成功！",
        ContinuousDays:  5,
        TotalDays:       15,
        RewardPoints:    10,
        CheckinTime:     time.Now().Unix(),
        IsFirstCheckin:  true,
    }, nil
}

func (s *checkinServer) GetCheckinStatus(ctx context.Context, req *pb.CheckinStatusRequest) (*pb.CheckinStatusResponse, error) {
    // TODO: 实现状态查询
    return &pb.CheckinStatusResponse{
        TodayCheckin:     false,
        ContinuousDays:   4,
        TotalDays:        14,
        LastCheckinTime:  time.Now().AddDate(0, 0, -1).Unix(),
        NextReward:       10,
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterCheckinServiceServer(s, &checkinServer{})
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

### 3. 客户端调用示例
```go
package main

import (
    "context"
    "log"
    "time"

    pb "bluebell_backend/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewCheckinServiceClient(conn)
    
    // 用户签到
    checkinResp, err := client.UserCheckin(context.Background(), &pb.UserCheckinRequest{
        UserId:    12345,
        DeviceId:  "device-123",
        IpAddress: "192.168.1.1",
        Timestamp: time.Now().Unix(),
        CheckinType: "normal",
    })
    if err != nil {
        log.Fatalf("checkin failed: %v", err)
    }
    log.Printf("签到结果: %+v", checkinResp)
    
    // 查询签到状态
    statusResp, err := client.GetCheckinStatus(context.Background(), &pb.CheckinStatusRequest{
        UserId: 12345,
    })
    if err != nil {
        log.Fatalf("get status failed: %v", err)
    }
    log.Printf("签到状态: %+v", statusResp)
}
```

## 📊 数据库设计建议

```sql
CREATE TABLE user_checkins (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    checkin_date DATE NOT NULL,
    continuous_days INT DEFAULT 1,
    reward_points INT DEFAULT 0,
    checkin_type VARCHAR(20) DEFAULT 'normal',
    device_id VARCHAR(100),
    ip_address VARCHAR(45),
    checkin_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_date (user_id, checkin_date),
    INDEX idx_user_id (user_id),
    INDEX idx_checkin_date (checkin_date)
);

CREATE TABLE user_checkin_stats (
    user_id BIGINT PRIMARY KEY,
    total_days INT DEFAULT 0,
    continuous_days INT DEFAULT 0,
    last_checkin_date DATE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 🎯 业务规则建议

### ✅ 签到规则
- 每日只能签到一次
- 连续签到中断后重新计算
- 不同连续天数对应不同奖励
- 支持补签（消耗积分或道具）

### 🏆 排行榜规则
- **连续签到榜** - 按continuous_days降序
- **总签到榜** - 按total_days降序  
- **月度榜** - 当月签到天数降序

### 🎁 奖励机制
```go
// 连续签到奖励配置
var rewardConfig = map[int32]int32{
    1:  5,   // 第1天
    3:  10,  // 第3天
    7:  20,  // 第7天
    15: 50,  // 第15天
    30: 100, // 第30天
}
```

## 🔄 后续扩展
- 签到日历展示
- 补签卡功能
- 签到分享奖励
- 特殊节日签到奖励翻倍
- 签到任务系统