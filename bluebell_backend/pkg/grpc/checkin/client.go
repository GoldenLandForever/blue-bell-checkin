package checkin

import (
	"bluebell_backend/pb"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	once     sync.Once
	instance *CheckinClient
)

type CheckinClient struct {
	conn   *grpc.ClientConn
	Client pb.CheckinServiceClient
}

// GetClient 获取签到的单例客户端
func InitClient(addr string) error {
	var err error
	once.Do(func() {
		conn, e := grpc.Dial(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(5*time.Second),
		)
		if e != nil {
			err = e
			return
		}

		instance = &CheckinClient{
			conn:   conn,
			Client: pb.NewCheckinServiceClient(conn),
		}
	})

	if err != nil {
		return err
	}
	return nil
}

//// UserCheckin 用户签到
//func (c *CheckinClient) UserCheckin(ctx context.Context, userID string) error {
//	//_, err := c.client.UserCheckin(ctx, &pb.UserCheckinRequest{
//	//	UserId: userID,
//	//})
//	return nil
//}

func GetClient() *CheckinClient {
	return instance
}

// Close 关闭连接
func (c *CheckinClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
