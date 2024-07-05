package handler

import (
	"html/template"
	"log"
	"net/http"
	"webchat/chat"
	"webchat/db"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type showLoginPageData struct {
	LoginInput        string
	ErrorLoginMessage string
}

func ShowLoginPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

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

	err = session.Save(request, response)
	if err != nil {
		log.Println(err)
	}

	err = loginPage.Execute(response, &showLoginPageData{LoginInput: loginInput, ErrorLoginMessage: errorLoginMessage})
	if err != nil {
		log.Println(err)
	}
}

type showChatPageData struct {
	UserName string
	IsBlock  bool
	IsAdmin  bool
}

func ShowChatPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {
	chatPage, err := template.ParseFiles("resources/chat.html")
	if err != nil {
		log.Println(err)
	}

	user := session.Values["user"].(db.User)

	err = chatPage.Execute(response, showChatPageData{
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

func Login(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

	passwordInput := request.FormValue("passwordInput")
	login := request.FormValue("loginInput")

	user, err := db.GetUserByLogin(login)

	if err != nil {
		log.Println(err)
	} else if checkPasswordHash(passwordInput, user.Password) {
		session.Values["authenticated"] = true
		session.Values["user"] = user
		err := session.Save(request, response)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(response, request, "/chat?command=show_chat_page", http.StatusFound)
		return
	}

	session.Values["LoginInput"] = login
	session.Values["ErrorLoginMessage"] = "Неправильный логин или пароль"

	err = session.Save(request, response)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(response, request, "/chat?command=show_login_page", http.StatusFound)
}

func checkPasswordHash(passwordInput, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordInput))
	return err == nil
}

func Logout(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

	session.Values["authenticated"] = false

	ShowLoginPage(response, request, session)
}
