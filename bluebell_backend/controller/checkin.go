package controller

import (
	"bluebell_backend/models"
	"bluebell_backend/pb"
	"bluebell_backend/pkg/grpc/checkin"
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
}
func CalendarHandler(c *gin.Context)    {}
func RetroactiveHandler(c *gin.Context) {}
