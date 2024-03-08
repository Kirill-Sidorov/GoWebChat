package main

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strings"
	"webchat/handler"
	"webchat/users"
)

var store *sessions.CookieStore

func main() {
	store = sessions.NewCookieStore([]byte("key-secret"))
	store.MaxAge(86400)

	gob.Register(users.User{})

	server := http.NewServeMux()
	server.HandleFunc("/", baseHandler)
	server.HandleFunc("/favicon.ico", faviconHandler)

	log.Println("Server Started...")
	err := http.ListenAndServe("localhost:8080", server)
	log.Fatal(err)
}

func baseHandler(writer http.ResponseWriter, request *http.Request) {

	command := request.URL.Query().Get("command")
	command = strings.TrimSpace(command)
	command = strings.ToUpper(command)

	session, _ := store.Get(request, "cookie-name")
	auth, ok := session.Values["authenticated"].(bool)

	if !ok || !auth {

		if command == "LOGIN" {
			handler.Login(writer, request, session)
			return
		}
		handler.ShowLoginPage(writer, request, session)
		return
	}

	switch command {
	case "SHOW_CHAT_PAGE":
		handler.ShowChatPage(writer, request, session)
	case "SEND_MESSAGE":
		handler.SendMessage(writer, request, session)
	case "LOGOUT":
		handler.Logout(writer, request, session)
	default:
		handler.ShowChatPage(writer, request, session)
	}
}

func faviconHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "resources/favicon.ico")
}
