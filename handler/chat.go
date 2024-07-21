package handler

import (
	"log"
	"net/http"
	"webchat/chat"
	"webchat/db"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type showChatPageData struct {
	UserName string
	IsBlock  bool
	IsAdmin  bool
}

func ShowChatPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

	user := session.Values["user"].(db.User)

	err := chatPageHtml.Execute(response, showChatPageData{
		UserName: user.Name,
		IsBlock:  false,
		IsAdmin:  user.Type == db.ADMIN,
	})

	if err != nil {
		log.Println(err)
	}
}

func CreateWebSocketConnection(response http.ResponseWriter,
	request *http.Request,
	session *sessions.Session) {

	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := session.Values["user"].(db.User)
	client := chat.NewClient(conn, user.Name)

	go client.WritePump()
	go client.ReadPump()
}