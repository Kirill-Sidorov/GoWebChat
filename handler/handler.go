package handler

import (
	"html/template"
	"log"
	"net/http"
	"webchat/chat"
	"webchat/users"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type showLoginPageData struct {
	LoginInput        string
	ErrorLoginMessage string
}

func ShowLoginPage(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {

	loginPage, err := template.ParseFiles("resources/login.html")
	if err != nil {
		log.Println(err)
	}

	loginInput, hasLogin := session.Values["LoginInput"].(string)
	errorLoginMessage, hasError := session.Values["ErrorLoginMessage"].(string)

	if !hasLogin {
		loginInput = ""
	}
	if !hasError {
		errorLoginMessage = ""
	}

	delete(session.Values, "LoginInput")
	delete(session.Values, "ErrorLoginMessage")

	err = session.Save(request, writer)
	if err != nil {
		log.Println(err)
	}

	err = loginPage.Execute(writer, showLoginPageData{LoginInput: loginInput, ErrorLoginMessage: errorLoginMessage})
	if err != nil {
		log.Println(err)
	}
}

type showChatPageData struct {
	UserName string
	IsBlock  bool
	IsAdmin  bool
}

func ShowChatPage(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {
	chatPage, err := template.ParseFiles("resources/chat.html")
	if err != nil {
		log.Println(err)
	}

	user := session.Values["user"].(users.User)

	err = chatPage.Execute(writer, showChatPageData{
		UserName: user.Name,
		IsBlock:  false,
		IsAdmin:  user.Type == users.ADMIN,
	})

	if err != nil {
		log.Println(err)
	}
}

/*
func SendMessage(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {
	user := session.Values["user"].(users.User)
	message := request.FormValue("message")
	message = strings.TrimSpace(message)
	if len(message) > 0 {
		messages.AddMessage(user.Name, message)
	}
	http.Redirect(writer, request, "/chat?command=show_chat_page", http.StatusFound)
}
*/

func CreateWebSocketConnection(writer http.ResponseWriter, 
							   request *http.Request, 
							   session *sessions.Session,
							   hub *chat.Hub) {

								
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := session.Values["user"].(users.User)
	client := chat.NewClient(hub, conn, user.Name)
	
	go client.WritePump()
	go client.ReadPump(hub)						
}

func Login(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {

	password := request.FormValue("passwordInput")
	login := request.FormValue("loginInput")

	user, find := users.ChatUsersMap[login]

	if find && user.Password == password {

		session.Values["authenticated"] = true
		session.Values["user"] = user
		err := session.Save(request, writer)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(writer, request, "/chat?command=show_chat_page", http.StatusFound)
		return
	}

	session.Values["LoginInput"] = login
	session.Values["ErrorLoginMessage"] = "Неправильный логин или пароль"

	err := session.Save(request, writer)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(writer, request, "/chat?command=show_login_page", http.StatusFound)
}

func Logout(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {

	session.Values["authenticated"] = false

	ShowLoginPage(writer, request, session)
}
