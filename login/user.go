package user

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/scrypt"
)

// ScryptN is the current recommended value for the Scrypt N parameter
const ScryptN = 16384

// ScryptN is the current recommended value for the Scrypt R parameter
const ScryptR = 8

// ScryptN is the current recommended value for the Scrypt P parameter
const ScryptP = 1

// ScryptN is the number of bytes to generate
const ScryptBytes = 512

// SaltBytes is the size of the salt to create
const SaltBytes = 16

// TODO don't overwrite existing
func CreateUser(db *sql.DB, name string, pass string) error {
	saltBts := make([]byte, SaltBytes)
	if _, err := rand.Read(saltBts); err != nil {
		return fmt.Errorf("failed to get salt: %v", err)
	}
	salt := saltBts
	// salt := base64.URLEncoding.EncodeToString(saltBts)

	hashedPass, err := scrypt.Key([]byte(pass), salt, ScryptN, ScryptR, ScryptP, ScryptBytes)
	if err != nil {
		return fmt.Errorf("failed to hash: %v", err)
	}

	if _, err := db.Exec("insert into user (name, salt, hash, updated) values($1, $2, $3, $4);", name, salt, hashedPass, time.Now().Unix()); err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

var BadPasswordErr = errors.New("bad password")

// Login attempts to log in, and returns a token on success. If the login fails, an error is returned. The error MUST NOT be passed to the user, for security.
func Login(db *sql.DB, tokenKey []byte, user string, pass string) (string, error) {
	if err := Authenticate(db, user, pass); err != nil {
		return "", BadPasswordErr
	}
	token, err := CreateToken(user, tokenKey)
	if err != nil {
		return "", fmt.Errorf("creating token; %v", err)
	}
	return token, nil
}

// Authenticate returns nil if the user has the given password in the database, else an error. This error is for debugging or logging and MUST NOT be passed to the user for sercurity.
func Authenticate(db *sql.DB, name string, pass string) error {
	q := "select salt, hash from user where name = $1;"
	salt := []byte{}
	dbHashedPass := []byte{}
	if err := db.QueryRow(q, name).Scan(&salt, &dbHashedPass); err != nil {
		return fmt.Errorf("querying: %v", err)
	}
	hashedPass, err := scrypt.Key([]byte(pass), salt, ScryptN, ScryptR, ScryptP, ScryptBytes)
	if err != nil {
		return fmt.Errorf("hash failure: %v", err)
	}
	if !bytes.Equal(hashedPass, dbHashedPass) {
		return fmt.Errorf("bad password")
	}
	return nil
}
