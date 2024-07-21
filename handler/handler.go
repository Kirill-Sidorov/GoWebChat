package handler

import (
	"html/template"
	"log"
)

var (
	loginPageHtml 		 *template.Template
	registrationPageHtml *template.Template
	chatPageHtml 		 *template.Template
)

func Init() {
	var err error
	loginPageHtml, err = template.ParseFiles("resources/login.html")
	if err != nil {
		log.Fatal(err)
	}

	registrationPageHtml, err = template.ParseFiles("resources/registration.html")
	if err != nil {
		log.Fatal(err)
	}

	chatPageHtml, err = template.ParseFiles("resources/chat.html")
	if err != nil {
		log.Fatal(err)
	}
}