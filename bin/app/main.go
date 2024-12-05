package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/louisfield/multi_backend/internal/services"
	"github.com/louisfield/multi_backend/internal/types"
	"github.com/lxzan/gws"
)

const (
	PingInterval         = 5 * time.Second
	HeartbeatWaitTimeout = 10 * time.Second
)

//go:embed index.html
var html []byte

func main() {
	var upgrader = gws.NewUpgrader(&WebSocket{}, &gws.ServerOption{
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},

		Authorize: func(r *http.Request, session gws.SessionStorage) bool {
			var name = r.URL.Query().Get("name")
			if name == "" {
				return false
			}
			var token = r.URL.Query().Get("token")

			if token == "" {
				return false
			}

			session.Store("name", name)

			session.Store("token", token)
			services.MaybeCreateLobby(session, token)
			return true
		},
	})

	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			log.Printf("Accept: " + err.Error())
			return
		}
		socket.ReadLoop()
	})

	http.HandleFunc("/index.html", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(html)
	})

	log.Printf("running")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf("running")
}

func MustLoad[T any](session gws.SessionStorage, key string) (v T) {
	if value, exist := session.Load(key); exist {
		v, _ = value.(T)
	}
	return
}

type WebSocket struct{}

func (c *WebSocket) OnOpen(socket *gws.Conn) {
	name := MustLoad[string](socket.Session(), "name")
	token := MustLoad[string](socket.Session(), "token")
	lobby := services.FindLobby(token)
	user := services.NewUser(name)
	user.Conn = socket

	lobby.Join(user)

	_ = socket.SetDeadline(time.Now().Add(PingInterval + HeartbeatWaitTimeout))

}

func (c *WebSocket) OnClose(socket *gws.Conn, err error) {
	name := MustLoad[string](socket.Session(), "name")

	log.Printf("onerror, name=%s, msg=%s\n", name, err.Error())
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + HeartbeatWaitTimeout))
	_ = socket.WriteString("pong")
}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	if b := message.Bytes(); len(b) == 4 && string(b) == "ping" {
		c.OnPing(socket, nil)
		return
	}

	var input = &types.Input{}

	json.Unmarshal(message.Bytes(), input)
	token := MustLoad[string](socket.Session(), "token")
	lobby := services.FindLobby(token)
	if lobby == nil {
		return
	}

	lobby.Broadcast(input.Message)

}
