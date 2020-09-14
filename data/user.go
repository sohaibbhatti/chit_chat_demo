package data

import "time"

type User struct {
	Id int
	Uuid string
	Name string
	Email string
        Password string
	CreatedAt time.Time
}

type Session struct {
	Id int
	Uuid string
	Email string
	UserId int
	CreatedAt time.Time
}


func (user *User) CreateSession() (session Session, err error) {
	statement := "INSERT INTO sessions (uuid, email, user_id, created_at) VALUES ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}


// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions where uuid=$1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions where uuid=$1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	if err != nil {
		valid = false
		return
	}

	if session.Id != 0 {
		valid = true
	}

	return
}

func (user *User) Create() (err error) {
	statement := "INSERT INTO users (uuid, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
