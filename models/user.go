package models

import (
	"awesomeProject3/db"
	"awesomeProject3/util"
	"errors"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (user User) Save() error {
	query := "INSERT INTO users(email, password) VALUES(?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(user.Email, hashPassword)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	user.ID = id
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)
	var retrievePassword string
	err := row.Scan(&u.ID, &retrievePassword)
	if err != nil {
		return err
	}
	passwordIsValid := util.CheckPasswordHash(u.Password, retrievePassword)
	if !passwordIsValid {
		return errors.New("Credentials is invalid")
	}

	return nil
}
