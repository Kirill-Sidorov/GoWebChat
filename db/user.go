package db

const (
	CLIENT = "CLIENT"
	ADMIN  = "ADMIN"
)

type User struct {
	Id		 int
	Login    string
	Password string
	Name     string
	Type     string
}

func GetUserByLogin(login string) (*User, error) {
	row := db.QueryRow("SELECT * FROM Client WHERE login = $1", login)
	user := User{}
	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Name, &user.Type)
	if err != nil{
		return nil, err
	}
	return &user, nil
}
