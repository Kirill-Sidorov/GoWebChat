package handler

import (
	"log"
	"net/http"
	"webchat/db"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type showRegistrationPage struct {
	NameInput			 string
	LoginInput 			 string
	ErrorRegisterMessage string
}

func ShowRegistrationPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {
	
	nameInput, hasName := session.Values["NameInput"].(string)
	loginInput, hasLogin := session.Values["LoginInput"].(string)
	errorRegisterMessage, hasError := session.Values["ErrorRegisterMessage"].(string)

	if !hasName {
		nameInput = ""
	}
	if !hasLogin {
		loginInput = ""
	}
	if !hasError {
		errorRegisterMessage = ""
	}

	delete(session.Values, "NameInput")
	delete(session.Values, "LoginInput")
	delete(session.Values, "ErrorRegisterMessage")

	err := session.Save(request, response)
	if err != nil {
		log.Println(err)
	}

	err = registrationPageHtml.Execute(response, &showRegistrationPage{
		NameInput: nameInput,
		LoginInput: loginInput, 
		ErrorRegisterMessage: errorRegisterMessage,
	})
	if err != nil {
		log.Println(err)
	}
}

func Register(response http.ResponseWriter, request *http.Request, session *sessions.Session) {

	nameInput := request.FormValue("nameInput")
	loginInput := request.FormValue("loginInput")
	passwordInput := request.FormValue("passwordInput")
	passwordRepeatInput := request.FormValue("passwordRepeatInput")

	session.Values["NameInput"] = nameInput
	session.Values["LoginInput"] = loginInput

	// empty checks
	if len(nameInput) == 0 {
		redirectWithErrorMessage(response, request, session, "Name input field empty")
		return
	}
	if len(loginInput) == 0 {
		redirectWithErrorMessage(response, request, session, "Login input field empty")
		return
	}
	if len(passwordInput) == 0 {
		redirectWithErrorMessage(response, request, session, "Password input field empty")
		return
	}
	if len(passwordRepeatInput) == 0 {
		redirectWithErrorMessage(response, request, session, "Password repeat input field empty")
		return
	}

	// size checks
	if len([]rune(nameInput)) < 4 {
		redirectWithErrorMessage(response, request, session, "Name must be more than 3 characters")
		return
	}
	if len([]rune(loginInput)) < 4 {
		redirectWithErrorMessage(response, request, session, "Login must be more than 3 characters")
		return
	}
	if len([]rune(passwordInput)) < 4 {
		redirectWithErrorMessage(response, request, session, "Password must be more than 3 characters")
		return
	}

	// already exists checks
	if db.IsExistUserWithName(nameInput) {
		redirectWithErrorMessage(response, request, session, "A user with the same name already exists, you must select a different name")
		return
	}
	if db.IsExistUserWithLogin(loginInput) {
		redirectWithErrorMessage(response, request, session, "This login is not suitable, you must select a new login")
		return
	}

	if passwordInput != passwordRepeatInput {
		redirectWithErrorMessage(response, request, session, "Passwords must match")
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordRepeatInput), 14)
	if err != nil {
		log.Println(err)
		redirectWithErrorMessage(response, request, session, "Failed to create new user")
		return
	}

	err = db.CreateNewUser(loginInput, string(bytes), nameInput)
	if err != nil {
		log.Println(err)
		redirectWithErrorMessage(response, request, session, "Failed to create new user")
		return
	}

	delete(session.Values, "NameInput")
	delete(session.Values, "LoginInput")
	delete(session.Values, "ErrorRegisterMessage")

	err = session.Save(request, response)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(response, request, "/chat?command=show_login_page", http.StatusFound)
}

func redirectWithErrorMessage(response http.ResponseWriter, request *http.Request, session *sessions.Session, message string) {
	session.Values["ErrorRegisterMessage"] = message
	err := session.Save(request, response)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(response, request, "/chat?command=show_registration_page", http.StatusFound)
}