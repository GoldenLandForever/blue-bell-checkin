package controller

import (
	"bluebell_backend/models"
	"bluebell_backend/pb"
	"bluebell_backend/pkg/grpc/checkin"
	"bluebell_backend/pkg/kafka"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DailyHandler(c *gin.Context) {
	instance := checkin.GetClient()
	var checkinmodel models.Checkin
	if err := c.ShouldBindJSON(&checkinmodel); err != nil {
		zap.L().Error("DailyHandler - ShouldBindJSON", zap.Error(err))
	}
	checkinResp, err := instance.Client.UserCheckin(c, &pb.UserCheckinRequest{UserId: checkinmodel.UserID, Timestamp: checkinmodel.TimeStamp, CheckinType: checkinmodel.CheckinType})
	if err != nil {
		zap.L().Error("could not checkin", zap.Error(err))
		return
	}

	zap.L().Info("签到响应详情",
		zap.Bool("success", checkinResp.Success),
		zap.String("message", checkinResp.Message),
		zap.Int32("continuous_days", checkinResp.ContinuousDays),
		zap.Int32("total_days", checkinResp.TotalDays),
		zap.Int32("remain_retro_times", checkinResp.RemainRetroTimes),
	)
	ResponseSuccess(c, nil)
	kafkaClient := kafka.GetKafkaClient()
	var checkinmodel2 models.CheckinResp
	checkinmodel2.UserID = checkinmodel.UserID
	checkinmodel2.ContinuousDays = checkinResp.ContinuousDays
	checkinmodel2.CheckinType = checkinmodel.CheckinType
	scoreData, err := json.Marshal(checkinmodel2)
	if err != nil {
		zap.L().Error("DailyHandler - json.Marshal", zap.Error(err))
	}
	if err := kafkaClient.SendMessage("Checkin", []byte(fmt.Sprintf("%d", checkinmodel2.UserID)), scoreData); err != nil {
		zap.L().Error("发送积分消息失败", zap.Error(err))
	}
}
func CalendarHandler(c *gin.Context)    {}
func RetroactiveHandler(c *gin.Context) {}
