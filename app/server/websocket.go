package server

import (
	"Meink/app/httprouter"
	"Meink/app/util"
	"net/http"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

/*
实现WebSockent ，减轻网络压力
*/
func Websocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if c, err := upgrader.Upgrade(w, r, nil); err != nil {
		util.MLog(err)
	} else {
		conn = c
	}

}
