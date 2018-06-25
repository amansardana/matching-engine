package ws

import (
	"github.com/amansardana/matching-engine/types"
	"github.com/gorilla/websocket"
)

type Ws struct {
	Conn        *websocket.Conn
	ReadChannel chan *types.WsMsg
}

var Connections map[string]*Ws
