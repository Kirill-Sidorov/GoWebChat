package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"webchat/chat"
	"webchat/handler"
	"webchat/users"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	initDb()

	store = sessions.NewCookieStore([]byte("key-secret"))
	store.MaxAge(86400)

	gob.Register(users.User{})

	hub = chat.NewHub()
	go hub.Run()

	server := http.NewServeMux()
	server.HandleFunc("/", baseHandler)
	server.HandleFunc("/favicon.ico", faviconHandler)

	log.Println("Server Started...")
	err = http.ListenAndServe("localhost:8080", server)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bye bye")
}

func initDb() *sql.DB {
	dbUser, found := os.LookupEnv("db_username")
	if !found {
		log.Panic("environment variable db_username not found in .env")
	}

	dbPassword, found := os.LookupEnv("db_password")
	if !found {
		log.Panic("environment variable db_password not found in .env")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=webchat sslmode=disable", dbUser, dbPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db
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
	case "CREATE_WEB_SOCKET_CONNECTION":
		handler.CreateWebSocketConnection(writer, request, session, hub)
	case "LOGOUT":
		handler.Logout(writer, request, session)
	default:
		handler.ShowChatPage(writer, request, session)
	}
}

func faviconHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "resources/favicon.ico")
}
