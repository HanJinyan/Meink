package app

import (
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

/*
实现WebSockent ，减轻网络压力
*/
func Websocket(ctx *Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if c, err := upgrader.Upgrade(ctx.Res, ctx.Req, nil); err != nil {
		MLog(err)
	} else {
		conn = c
	}
	ctx.Stop()
}
