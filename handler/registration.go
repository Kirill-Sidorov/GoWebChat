package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func ShowRegistrationPage(response http.ResponseWriter, request *http.Request, session *sessions.Session) {
	
	var data struct{}
	err := registrationPageHtml.Execute(response, data)

	if err != nil {
		log.Println(err)
	}
}