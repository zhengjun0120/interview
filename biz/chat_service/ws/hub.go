package ws

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"
	"ai_interview/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	HR        string = "hr"
	CANDIDATE string = "candidate" //候选者

	CHAT  string = "chat"  //聊天
	ERROR string = "error" //错误
	JOIN  string = "join"  //加入
	LEAVE string = "leave" //离开
)

type ChatMsgData struct {
	RoomID string `json:"room_id"`
	Text   string `json:"text"`
}

type Message struct {
	Type      string      `json:"type"`       //消息类型
	UserID    string      `json:"user_id"`    //用户ID
	MessageID string      `json:"message_id"` //消息ID
	Data      ChatMsgData `json:"data"`       //消息数据
	From      string      `json:"from"`       //消息来源
	Error     string      `json:"error"`      //错误信息
	ErrorCode int         `json:"error_code"` //错误代码
}

type Client struct {
	RoomID string          //房间ID
	UserID string          //用户ID
	Role   string          //用户角色
	Conn   *websocket.Conn //连接
	Send   chan []byte     //发送通道 将消息发送到该连接
}

type Room struct { // 房间 只允许有俩个连接
	Hr        *Client
	Candidate *Client
}

type Hub struct {
	Rooms      map[string]*Room   // 房间 map 键为房间ID 值为房间
	Clients    map[string]*Client // 连接 map 键为用户ID 值为连接
	Register   chan *Client       // 注册通道 新连接通过该通道注册
	Unregister chan *Client       // 注销通道 连接通过该通道注销
	Broadcast  chan Message       // 广播通道 用于将消息广播给所有连接
	Mu         sync.RWMutex       // 读写锁 保护 Rooms 和 Clients

	ChatRepo repo.ChatRepo
}

var hub *Hub

func InitWebSocketHub(chatRepo repo.ChatRepo) {
	hub = &Hub{
		Rooms:      make(map[string]*Room),
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		Mu:         sync.RWMutex{},
		ChatRepo:   chatRepo,
	}

	zlog.Infof("WebSocket hub 初始化成功")
}

func GetWebSocketHub() *Hub {
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case msg := <-h.Broadcast:
			h.broadcastToRoom(msg)
		}
	}
}

func (h *Hub) registerClient(c *Client) {
	h.Mu.Lock()

	// 检查房间是否存在
	_, ok := h.Rooms[c.RoomID]
	if !ok {
		h.Rooms[c.RoomID] = &Room{}
	}

	room := h.Rooms[c.RoomID]
	var alreadyConnected bool

	// 检查连接是否存在
	if c.Role == HR {
		if room.Hr != nil {
			alreadyConnected = true
		} else {
			room.Hr = c
		}
	} else if c.Role == CANDIDATE {
		if room.Candidate != nil {
			alreadyConnected = true
		} else {
			room.Candidate = c
		}
	}

	// 如果连接已存在，则关闭连接
	if alreadyConnected {
		h.Mu.Unlock()

		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "您已在其它设备连接")
		_ = c.Conn.WriteMessage(websocket.CloseMessage, closeMsg)
		_ = c.Conn.Close()
		close(c.Send)
		return
	}

	h.Clients[c.UserID] = c

	go c.ReadPump(h)
	go c.WritePump()

	var targets []*Client
	if room.Hr != nil {
		targets = append(targets, room.Hr)
	}
	if room.Candidate != nil {
		targets = append(targets, room.Candidate)
	}

	h.Mu.Unlock()

	// 构建加入消息
	joinMsg := Message{
		Type: JOIN,
		From: c.Role,
	}

	// 发送加入消息给房间内的俩个连接
	joinMsgBytes, err := json.Marshal(joinMsg)
	if err != nil {
		zlog.Errorf("解析消息错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(joinMsgBytes))
		return
	}
	fmt.Println("hr是否在线", room.Hr != nil)
	fmt.Println("candidate是否在线", room.Candidate != nil)
	for _, target := range targets {
		select {
		case target.Send <- joinMsgBytes:
		default:
			zlog.Errorf("发送消息错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(joinMsgBytes))
		}
	}
}

func (h *Hub) unregisterClient(c *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	room, ok := h.Rooms[c.RoomID]
	if ok {
		// 构建离开消息
		leaveMsg, _ := json.Marshal(Message{
			Type:   LEAVE,
			UserID: c.UserID,
			From:   c.Role,
		})

		if c.Role == HR {
			// 如果HR离开，则通知Candidate，如果Candidate存在的话
			if room.Candidate != nil {
				room.Candidate.Send <- leaveMsg
			}
			close(room.Hr.Send)

			room.Hr = nil
		} else if c.Role == CANDIDATE {
			// 如果Candidate离开，则通知HR，如果HR存在的话
			if room.Hr != nil {
				room.Hr.Send <- leaveMsg
			}
			close(room.Candidate.Send)
			room.Candidate = nil
		}

		if room.Hr == nil && room.Candidate == nil {
			delete(h.Rooms, c.RoomID)
		}
	}

	// 删除连接
	_, ok = h.Clients[c.UserID]
	if ok {
		delete(h.Clients, c.UserID)
	}

}

func (h *Hub) broadcastToRoom(msg Message) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	client, _ := h.Clients[msg.UserID]

	room, ok := h.Rooms[client.RoomID]
	if !ok {

		client.SendErrorMsg(response.ROOM_NOT_EXIST)
		return
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		client.SendErrorMsg(response.PARSE_MESSAGE_ERROR)
		zlog.Errorf("解析消息错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, msg.UserID, msg.Data.RoomID, string(msgBytes))
		return
	}

	// 保存聊天记录 入库
	err = h.ChatRepo.SaveInterviewMessage(context.Background(), entity.InterviewMessage{
		UserID:    client.UserID,
		RoomID:    client.RoomID,
		From:      client.Role,
		MessageID: msg.MessageID,
		Type:      msg.Type,
		Text:      msg.Data.Text,
	})
	if err != nil {
		client.SendErrorMsg(response.CustomError(fmt.Errorf("消息发送失败, %v", err)))
		zlog.Errorf("保存聊天记录错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, msg.UserID, msg.Data.RoomID, string(msgBytes))
		return
	}

	// 广播消息给房间内的俩个连接
	if msg.From == HR && room.Candidate != nil {
		room.Candidate.Send <- msgBytes
	} else if msg.From == CANDIDATE && room.Hr != nil {
		room.Hr.Send <- msgBytes
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zlog.Errorf("读取错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(message))
			}
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			c.SendErrorMsg(response.PARSE_MESSAGE_ERROR)
			zlog.Errorf("解析消息错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(message))
			continue
		}

		if msg.MessageID == "" {
			msg.MessageID = util.GenerateStringID()
		}
		msg.Type = CHAT

		// 检查用户ID是否已连接
		userID := c.UserID
		if userID == "" {
			c.SendErrorMsg(response.USER_ID_IS_EMPTY)
			zlog.Errorf("用户ID为空 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(message))
			continue
		} else {
			_, ok := hub.Clients[userID]
			if !ok {
				c.SendErrorMsg(response.USER_NOT_CONNECTED)
				zlog.Errorf("用户未连接 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(message))
				continue
			}
		}

		msg.UserID = userID
		msg.From = c.Role

		hub.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Conn.Close()
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "写入消息时错误（服务端错误）")
				_ = c.Conn.WriteMessage(websocket.CloseMessage, closeMsg)

				zlog.Errorf("写入错误 error: %v  用户ID：%s  房间ID：%s  原始消息：%s", err, c.UserID, c.RoomID, string(message))
				return
			}
		}
	}
}

func (c *Client) SendErrorMsg(msgCode response.MsgCode) {
	errorMsg := Message{
		Type:      ERROR,
		Error:     msgCode.Msg,
		ErrorCode: msgCode.Code,
	}
	errorMsgBytes, _ := json.Marshal(errorMsg)
	c.Send <- errorMsgBytes
}
