package main

const (
	CLIENT uint8 = iota
	ADMIN
)

type User struct {
	Login    string
	Password string
	Type     uint8
	Name     string
}

var users = map[string]User{
	"admin": User{Login: "admin", Password: "admin", Type: ADMIN, Name: "Админ"},
	"anton": User{Login: "anton", Password: "anton", Type: CLIENT, Name: "Антон"},
}
