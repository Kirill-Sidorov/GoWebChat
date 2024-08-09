package db

const (
	CLIENT = "CLIENT"
	ADMIN  = "ADMIN"
)

type User struct {
	Id       int
	Login    string
	Password string
	Name     string
	Type     string
}

func GetUserByLogin(login string) (*User, error) {
	row := db.QueryRow("SELECT * FROM Client WHERE login = $1", login)
	user := User{}
	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Name, &user.Type)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsExistUserWithLogin(login string) bool {
	row := db.QueryRow("SELECT login FROM Client WHERE login = $1", login)
	var foundLogin string
	err := row.Scan(&foundLogin)
	return err == nil
}

func IsExistUserWithName(name string) bool {
	row := db.QueryRow("SELECT name FROM Client WHERE name = $1", name)
	var foundName string
	err := row.Scan(&foundName)
	return err == nil
}

func CreateNewUser(login, password, name string) error {
	_, err := db.Exec("INSERT INTO Client (login, password, name, type) VALUES ($1, $2, $3, $4)",
		login, password, name, CLIENT)
	return err
}