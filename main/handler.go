package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
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

	err = loginPage.Execute(writer, ShowLoginPageData{LoginInput: loginInput, ErrorLoginMessage: errorLoginMessage})
	if err != nil {
		log.Println(err)
	}

	delete(session.Values, "LoginInput")
	delete(session.Values, "ErrorLoginMessage")

	err = session.Save(request, writer)
	if err != nil {
		log.Println(err)
	}
}

func ShowChatPage(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("chat page"))
}

func Login(writer http.ResponseWriter, request *http.Request, session *sessions.Session) {

	password := request.FormValue("passwordInput")
	login := request.FormValue("loginInput")

	user, find := users[login]

	if find && user.Password == password {

		session.Values["authenticated"] = true
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
