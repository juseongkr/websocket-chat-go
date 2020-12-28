package db

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("DB: Unauthorized")

func CreateUser(username, password string) (id int, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	err = db.QueryRow(
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		username, string(hash),
	).Scan(&id)

	return
}

func FindUser(username, password string) (id int, err error) {
	var hash string
	err = db.QueryRow(
		`SELECT id, password FROM users WHERE username = $1`,
		username,
	).Scan(&id, &hash)

	if err == sql.ErrNoRows {
		return -1, ErrUnauthorized
	} else if err != nil {
		return -1, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return -1, ErrUnauthorized
	}

	return
}
