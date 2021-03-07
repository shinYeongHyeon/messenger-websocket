package db

import (
	"database/sql"
	"errors"
)

// CreateUser 새 사용자 만듦
func CreateUser(username, password string) (id int, err error) {
	err = db.QueryRow(`
	INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
	username,
	password).Scan(&id)

	return
}

// ErrUnauthorized 권한 없음 에러
var ErrUnauthorized = errors.New("db: unauthorized")

// FindUser 닉네임, 비번 가지고 사용자 찾기
func FindUser(username, password string) (id int, err error) {
	err = db.QueryRow(`
	SELECT id, password FROM users WHERE username = $1 and password = $2`,
	username,
	password).Scan(&id, &password)

	if err == sql.ErrNoRows {
		return 0, ErrUnauthorized
	} else if err != nil {
		return 0, err
	}

	return
}

