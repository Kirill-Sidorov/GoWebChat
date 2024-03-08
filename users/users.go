package users

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

var ChatUsersMap = map[string]User{
	"admin": {Login: "admin", Password: "admin", Type: ADMIN, Name: "Админ"},
	"anton": {Login: "anton", Password: "anton", Type: CLIENT, Name: "Антон"},
}
