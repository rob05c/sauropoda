package user

import (
	"crypto/rand"
	"database/sql"
	"fmt"
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

func CreateUser(db *sql.DB, name string, pass string) error {
	salt := make([]byte, SaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to get salt: %v", err)
	}
	hashedPass, err := scrypt.Key([]byte(pass), salt, ScryptN, ScryptR, ScryptP, ScryptBytes)
	if err != nil {
		return fmt.Errorf("failed to hash: %v", err)
	}

	if _, err := db.Exec("insert into user (name, salt, hash) values($1, $2, $3);", name, salt, hashedPass); err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}
