package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
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

}