package main

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"webchat/chat"
	"webchat/db"
	"webchat/handler"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const (
	listenAddr      = ":8080"
	shutdownTimeout = 3 * time.Second
)

var store *sessions.CookieStore

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	handler.Init()
	db.Init()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runWebServer(ctx); err != nil {
		log.Println(err)
	}

	db.Close()
}

func runWebServer(ctx context.Context) error {

	secret, found := os.LookupEnv("cookie_store_secret")
	if !found {
		return errors.New("environment variable cookie_store_secret not found in .env")
	}
	store = sessions.NewCookieStore([]byte(secret))
	store.MaxAge(86400)

	gob.Register(db.User{})

	go chat.Run()

	var mux = http.NewServeMux()
	var srv = &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	mux.HandleFunc("/", baseHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			db.Close()
			log.Fatalf("Server Listen and Serve: %v", err)
		}
	}()
	log.Printf("Server Started: %s", listenAddr)

	<-ctx.Done()
	log.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	<-shutdownCtx.Done()

	err := ctx.Err()
	if err != nil {
		return err
	}

	return nil
}

func baseHandler(response http.ResponseWriter, request *http.Request) {

	command := request.URL.Query().Get("command")
	command = strings.TrimSpace(command)
	command = strings.ToUpper(command)

	session, _ := store.Get(request, "cookie-name")
	auth, ok := session.Values["authenticated"].(bool)

	if !ok || !auth {

		switch command {
		case "LOGIN":
			handler.Login(response, request, session)
		case "SHOW_REGISTRATION_PAGE":
			handler.ShowRegistrationPage(response, request, session)
		case "REGISTER":
			handler.Register(response, request, session)
		case "SHOW_LOGIN_PAGE":
			handler.ShowLoginPage(response, request, session)
		default:
			handler.ShowLoginPage(response, request, session)
		}
		return
	}

	switch command {
	case "SHOW_CHAT_PAGE":
		handler.ShowChatPage(response, request, session)	
	case "CREATE_WEB_SOCKET_CONNECTION":
		handler.CreateWebSocketConnection(response, request, session)
	case "LOGOUT":
		handler.Logout(response, request, session)
	default:
		handler.ShowChatPage(response, request, session)
	}
}

func faviconHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "resources/favicon.ico")
}
