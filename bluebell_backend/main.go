package main

import (
	"bluebell_backend/controller"
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/logger"
	"bluebell_backend/pkg/kafka"
	"bluebell_backend/pkg/snowflake"
	"bluebell_backend/routers"
	"bluebell_backend/settings"
	"fmt"

	"go.uber.org/zap"
)

// @host 127.0.0.1:8081
// @BasePath /api/v1/
func main() {

	//var confFile string
	//flag.StringVar(&confFile, "conf", "./conf/config.yaml", "配置文件")
	//flag.Parse()
	// 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	defer redis.Close()

	// Redis预热
	//warmUp := redis.NewRedisWarmUp()
	//if err := warmUp.WarmUpAll(); err != nil {
	//	zap.L().Error("Redis预热失败", zap.Error(err))
	//	// 预热失败不中断启动，只记录错误
	//}

	// 初始化定时任务
	if err := redis.SyncExpiringVotes(7 * 24 * 3600); err != nil {
		zap.L().Error("同步失败", zap.Error(err))
	} else {
		zap.L().Info("同步成功")
	}

	// 程序退出时关闭Kafka客户端
	defer kafka.CloseKafka()
	// 雪花算法生成分布式ID
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	//kafka初始化
	if err := kafka.InitKafka(); err != nil {
		fmt.Printf("init kafka failed, err: %v\n", err)
		return
	}
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator Trans failed,err:%v\n", err)
		return
	}
	// 注册路由
	r := routers.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
