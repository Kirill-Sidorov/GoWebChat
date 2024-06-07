package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"strings"
	"webchat/chat"
	"webchat/db"
	"webchat/handler"
	"webchat/users"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	store *sessions.CookieStore
	hub   *chat.Hub
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(".env file not found")
	}

	db.Init()

	secret, found := os.LookupEnv("cookie_store_secret")
	if !found {
		log.Panic("environment variable cookie_store_secret not found in .env")
	}
	store = sessions.NewCookieStore([]byte(secret))
	store.MaxAge(86400)

	gob.Register(users.User{})

	server := http.NewServeMux()
	server.HandleFunc("/", baseHandler)
	server.HandleFunc("/favicon.ico", faviconHandler)

	hub = chat.NewHub()
	go hub.Run()

	log.Println("Server Started: http://localhost:8080")
	err = http.ListenAndServe("localhost:8080", server)
	if err != nil {
		log.Fatal(err)
	}
}

func baseHandler(response http.ResponseWriter, request *http.Request) {

	command := request.URL.Query().Get("command")
	command = strings.TrimSpace(command)
	command = strings.ToUpper(command)

	session, _ := store.Get(request, "cookie-name")
	auth, ok := session.Values["authenticated"].(bool)

	if !ok || !auth {

		if command == "LOGIN" {
			handler.Login(response, request, session)
			return
		}
		handler.ShowLoginPage(response, request, session)
		return
	}

	switch command {
	case "SHOW_CHAT_PAGE":
		handler.ShowChatPage(response, request, session)
	case "CREATE_WEB_SOCKET_CONNECTION":
		handler.CreateWebSocketConnection(response, request, session, hub)
	case "LOGOUT":
		handler.Logout(response, request, session)
	default:
		handler.ShowChatPage(response, request, session)
	}
}

func faviconHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "resources/favicon.ico")
}
