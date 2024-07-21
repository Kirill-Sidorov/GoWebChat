package handler

import (
	"log"
	"net/http"
	"webchat/db"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type showLoginPageData struct {
	LoginInput        string
	ErrorLoginMessage string
}

func ShowLoginPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

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

	err := session.Save(request, response)
	if err != nil {
		log.Println(err)
	}

	err = loginPageHtml.Execute(response, &showLoginPageData{LoginInput: loginInput, ErrorLoginMessage: errorLoginMessage})
	if err != nil {
		log.Println(err)
	}
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
	session.Values["ErrorLoginMessage"] = "Incorrect login or password"

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