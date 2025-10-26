package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 表示一个WebSocket客户端连接
type Client struct {
	ID     string
	UserID uint64
	Socket *websocket.Conn
	Send   chan []byte
}

// Message WebSocket消息结构
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 心跳超时时间
const (
	writeWait      = 10 * time.Second    // 写超时
	pongWait       = 60 * time.Second    // 读超时
	pingPeriod     = (pongWait * 9) / 10 // 发送心跳间隔
	maxMessageSize = 512                 // 最大消息大小
)

// WebSocketHandler 处理WebSocket连接
func WebSocketHandler(c *gin.Context) {
	// 获取当前用户ID
	uidStr, exists := c.Get("userID")
	if !exists {
		ResponseError(c, CodeNotLogin)
		return
	}
	userID := uidStr.(uint64)

	// 升级HTTP连接为WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("upgrade websocket failed", zap.Error(err))
		return
	}

	// 创建新的客户端
	client := &Client{
		ID:     generateClientID(), // 这里可以使用UUID或其他方式生成唯一ID
		UserID: userID,
		Socket: ws,
		Send:   make(chan []byte, 256),
	}

	// 注册客户端到全局管理器
	manager.RegisterClient(client)

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// readPump 处理从WebSocket读取数据
func (c *Client) readPump() {
	defer func() {
		manager.UnregisterClient(c)
		c.Socket.Close()
	}()

	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Error("websocket read error", zap.Error(err))
			}
			break
		}

		// 处理接收到的消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			zap.L().Error("unmarshal message failed", zap.Error(err))
			continue
		}

		// 根据消息类型处理
		switch msg.Type {
		case "ping":
			// 收到ping消息，回复pong
			response := Message{
				Type: "pong",
				Data: time.Now().Unix(),
			}
			data, _ := json.Marshal(response)
			c.Send <- data
		}
	}
}

// writePump 处理向WebSocket写入数据
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
