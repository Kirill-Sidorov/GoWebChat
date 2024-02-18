package main

type User struct {
	Login    string
	Password string
	Name     string
}

var users = map[string]User{
	"admin": User{Login: "admin", Password: "admin", Name: "Админ"},
	"anton": User{Login: "anton", Password: "anton", Name: "Антон"},
}
