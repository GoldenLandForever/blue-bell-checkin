package controller

import (
	"sync"
)

// ClientManager 管理所有WebSocket连接
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.RWMutex
}

var manager = NewManager()

// GetManager 获取WebSocket管理器实例
func GetManager() *ClientManager {
	return manager
}

// NewManager 创建新的连接管理器
func NewManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// RegisterClient 注册新的客户端连接
func (manager *ClientManager) RegisterClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.Clients[client] = true
	manager.Register <- client
}

// UnregisterClient 注销客户端连接
func (manager *ClientManager) UnregisterClient(client *Client) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if _, ok := manager.Clients[client]; ok {
		manager.Unregister <- client
	}
}

// GetClientByUserID 根据用户ID获取客户端连接
func (manager *ClientManager) GetClientByUserID(userID uint64) []*Client {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	var clients []*Client
	for client := range manager.Clients {
		if client.UserID == userID {
			clients = append(clients, client)
		}
	}
	return clients
}

// Broadcast 广播消息给所有客户端
func (manager *ClientManager) BroadcastMessage(message []byte) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	for client := range manager.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(manager.Clients, client)
		}
	}
}

// SendToUser 发送消息给指定用户
func (manager *ClientManager) SendToUser(userID uint64, message []byte) {
	clients := manager.GetClientByUserID(userID)
	for _, client := range clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			manager.UnregisterClient(client)
		}
	}
}

// Run 启动WebSocket管理器
func (manager *ClientManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client] = true
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			for client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.Clients, client)
				}
			}
		}
	}
}
