package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sofferjacob/maker_api/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int          `db:"id" json:"id"`
	Email     string       `db:"email" json:"email"`
	Name      string       `db:"name" json:"name"`
	Password  string       `db:"password"`
	Joined    time.Time    `db:"joined" json:"joined"`
	LastLogin sql.NullTime `db:"last_login" json:"last_login"`
	Ts        string       `db:"ts"`
}

type UserData struct {
	Id        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Joined    time.Time `db:"joined" json:"joined"`
	LastLogin time.Time `db:"last_login" json:"last_login"`
}

func (u *User) ToUserData() UserData {
	return UserData{
		u.Id,
		u.Email,
		u.Name,
		u.Joined,
		u.LastLogin.Time,
	}
}

// Takes a struct with id set to the uid to update,
// and optional update values name, email
func (u *User) Update() error {
	if u.Name == "" && u.Email == "" {
		return nil
	}
	updateI := 1
	addComma := false
	query := "UPDATE users SET "
	args := []interface{}{}
	if u.Name != "" {
		query += fmt.Sprintf("name = $%v", updateI)
		updateI++
		addComma = true
		args = append(args, u.Name)
	}
	if u.Email != "" {
		if addComma {
			query += ", "
		}
		query += fmt.Sprintf("email = $%v", updateI)
		updateI++
		args = append(args, u.Email)

	}
	query += fmt.Sprintf(" WHERE id = $%v", updateI)
	args = append(args, u.Id)
	_, err := db.Client.Client.Exec(query, args...)
	return err
}

func (u *User) Register() error {
	if u.Email == "" || u.Password == "" || u.Name == "" {
		return errors.New("missing required struct fields (name, email, password)")
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err.Error())
	}
	_, err = db.Client.Client.Exec(
		"INSERT INTO users  (email, name, password) VALUES ($1, $2, $3);",
		u.Email,
		u.Name,
		string(pwd),
	)
	return err
}

func (u *User) Login() (string, error) {
	if u.Email == "" || u.Password == "" {
		return "", errors.New("missing required struct fields (email, password)")
	}
	pwd := []byte(u.Password)
	err := db.Client.Client.Get(u, "SELECT * FROM users WHERE email = $1;", u.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), pwd); err != nil {
		return "", fmt.Errorf("invalid password: %v", err.Error())
	}
	claims := jwt.StandardClaims{
		Issuer:    "a01028653@tec.mx",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(48 * time.Hour)).Unix(),
		Subject:   strconv.Itoa(u.Id),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(os.Getenv("AUTH_KEY")))
	if err != nil {
		return "", fmt.Errorf("could not issue token: %v", err.Error())
	}
	_, err = db.Client.Client.Exec("UPDATE users SET last_login = $1 WHERE id = $2;", time.Now(), u.Id)
	if err != nil {
		return token, fmt.Errorf("token created, failed to update last login: %v", err.Error())
	}
	return token, nil
}

func GetUser(id string) (UserData, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return UserData{}, err
	}
	u := User{}
	err = db.Client.Client.Get(&u, "SELECT * FROM users WHERE id = $1;", uid)
	return u.ToUserData(), err
}

func QueryUserFTS(query string) ([]UserData, error) {
	u := []User{}
	err := db.Client.Client.Select(&u, "SELECT * FROM query_gin(null::users, $1);", query)
	if err != nil {
		return nil, err
	}
	res := make([]UserData, 0, len(u))
	for _, v := range u {
		res = append(res, v.ToUserData())
	}
	return res, nil
}
