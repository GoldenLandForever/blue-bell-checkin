package controller

import (
	"fmt"
	"sync/atomic"
	"time"
)

var clientIDCounter uint64

// generateClientID 生成唯一的客户端ID
func generateClientID() string {
	timestamp := time.Now().UnixNano() / 1000000 // 毫秒时间戳
	count := atomic.AddUint64(&clientIDCounter, 1)
	return fmt.Sprintf("%d-%d", timestamp, count)
}
