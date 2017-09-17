package login

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const GUIDBytes = 16

var SigningMethod = jwt.SigningMethodHS512

const CookieName = "token"

const ExpirationLength = time.Hour * 24

func CreateToken(user string, key []byte) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(SigningMethod, &jwt.StandardClaims{
		Subject:   user,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(ExpirationLength).Unix(),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("creating jwt string: %v", err)
	}
	return tokenStr, nil
}

func CreateCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(ExpirationLength),
		HttpOnly: true,
		Path:     "/",
	}
}

// TokenValid returns the username if the token is valid.
func TokenValid(tokenStr string, db *sql.DB, key []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if signingMethod, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || signingMethod.Name != SigningMethod.Name {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return "", fmt.Errorf("parsing token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token invalid")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("unknown claims %T", token.Claims)
	}
	if err := claims.Valid(); err != nil {
		return "", fmt.Errorf("claims invalid: %v", err)
	}
	if err := dbValid(db, claims.Subject, claims.IssuedAt); err != nil {
		return "", fmt.Errorf("invalid in database: %v", err)
	}
	return claims.Subject, nil
}

// dbValid returns whether the given user exists in the database, and their password hasn't change since the given time.
func dbValid(db *sql.DB, name string, notModifiedAfter int64) error {
	q := "select updated from user where name = $1;"
	updated := int64(0)
	if err := db.QueryRow(q, name).Scan(&updated); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no such user")
		}
		return fmt.Errorf("querying: %v", err)
	}
	if updated > notModifiedAfter {
		return fmt.Errorf("password has changed")
	}
	return nil
}
