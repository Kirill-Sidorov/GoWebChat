package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"webchat/messages"
)

type ShowLoginPageData struct {
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

	err = loginPage.Execute(writer, ShowLoginPageData{LoginInput: loginInput, ErrorLoginMessage: errorLoginMessage})
	if err != nil {
		log.Println(err)
	}
}

type ShowChatPageData struct {
	UserName string
	Messages string
	IsBlock  bool
	IsAdmin  bool
}

func ShowChatPage(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {
	chatPage, err := template.ParseFiles("resources/chat.html")
	if err != nil {
		log.Println(err)
	}

	user := session.Values["user"].(User)

	err = chatPage.Execute(writer, ShowChatPageData{
		UserName: user.Name,
		Messages: messages.GetMessages(),
		IsBlock:  false,
		IsAdmin:  user.Type == ADMIN,
	})

	if err != nil {
		log.Println(err)
	}
}

func Login(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {

	password := request.FormValue("passwordInput")
	login := request.FormValue("loginInput")

	user, find := users[login]

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
