package main

import (
	"fmt"
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/serialize/json"
	"github.com/lonng/nano/session"
	"log"
	"net/http"
	"strings"
	"time"
)
type (
	Room struct {
		group *nano.Group
		component.Base
	}
	RoomManager struct {
		component.Base
		timer *scheduler.Timer
		rooms map[int]*Room
	}
	UserMessage struct {
		Name string `json:"name"`
		Content string `json:"content"`
	}

	NewUser struct {
		Content string `json:"content"`
	}

	AllMembers struct {
		Members []int64 `json:"members"`
	}

	JoinResponse struct {
		Code int `json:"code"`
		Result string `json:"result"`
	}

	stats struct {
		component.Base
		timer *scheduler.Timer
		outboundBytes int
		inboundBytes int
	}
)

func (stats *stats) outbound(s *session.Session, msg *pipeline.Message) error  {
	stats.outboundBytes += len(msg.Data)
	return nil
}

func (stats *stats) inbound (s *session.Session, msg *pipeline.Message) error {
	stats.inboundBytes += len(msg.Data)
	return nil
}

func (stats *stats) AfterInit()  {
	scheduler.NewTimer(time.Minute, func() {
		println("OutBoundBytes", stats.outboundBytes)
		println("InboundBytes", stats.inboundBytes)
	})
}

const (
	testRoomID = 1
	roomIDKey = "ROOM_ID"
)

// NewRoom 创建房间
func NewRoom() *Room {
	return &Room{
		group: nano.NewGroup("room"),
	}
}
// NewRoomManager 创建房间管理
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: map[int]*Room{},
	}
}
// AfterInit
func (mgr *RoomManager) AfterInit()  {
	session.Lifetime.OnClosed(func(s *session.Session) {
		if !s.HasKey(roomIDKey) {
			return
		}
		room := s.Value(roomIDKey).(*Room)
		err := room.group.Leave(s)
		if err != nil {
			panic(err)
		}
	})
	mgr.timer = scheduler.NewTimer(time.Minute, func() {
		for roomId, room := range mgr.rooms {
			println(fmt.Sprintf("UserCount: RoomID=%d, Time=%s,Count=%d", roomId, time.Now().String(), room.group.Count()))
		}
	})
}
// Join 加入
func (mgr *RoomManager) Join(s *session.Session, msg []byte) error {
	room, found := mgr.rooms[testRoomID]
	if !found {
		room = &Room{
			group: nano.NewGroup(fmt.Sprintf("room-%d", testRoomID)),
		}
		mgr.rooms[testRoomID] = room
	}
	fakeUID := s.ID()      //s.ID 用于uid
	_ = s.Bind(fakeUID) // 绑定id到room
	s.Set(roomIDKey, room)
	_ = s.Push("onMembers", &AllMembers{Members: room.group.Members()})
	//通知其他用户
	_ = room.group.Broadcast("onNewUser", &NewUser{Content: fmt.Sprintf("New user: %d", s.ID())})
	//新用户入组
	_ = room.group.Add(s)
	return s.Response(&JoinResponse{Result: "success"})
}

func (mgr *RoomManager) Message(s *session.Session, msg *UserMessage) error {
	if !s.HasKey(roomIDKey) {
		return fmt.Errorf("not join room yet")
	}
	room := s.Value(roomIDKey).(*Room)
	return room.group.Broadcast("onMessage", msg)
}

func main()  {
	components := &component.Components{}
	components.Register(NewRoomManager(), component.WithName("room"), component.WithNameFunc(strings.ToLower))
	pip := pipeline.New()
	var stats = &stats{}
	pip.Outbound().PushBack(stats.outbound)
	pip.Inbound().PushBack(stats.inbound)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	nano.Listen(":3250",
		nano.WithIsWebsocket(true),
		nano.WithPipeline(pip),
		nano.WithCheckOriginFunc(func(request *http.Request) bool {return true}),
		nano.WithWSPath("/nano"),
		nano.WithDebugMode(),
		nano.WithSerializer(json.NewSerializer()),
		nano.WithComponents(components),
	)
}
