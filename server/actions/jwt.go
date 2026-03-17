package actions

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginClaims struct {
	jwt.RegisteredClaims
}

type TokenPayload struct {
	Username string `json:"username"`
}

type Token struct {
	Token   string    `json:"token"`
	Expiry  int       `json:"expiry"`
	Created time.Time `json:"created"`
}

func GenerateJWTToken(data TokenPayload, expiry int) (*Token, error) {
	now := time.Now().UTC()
	registeredClaims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(now),
		Subject:  data.Username,
	}

	if expiry != -1 {
		registeredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Second * time.Duration(expiry)))
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, LoginClaims{
		registeredClaims,
	})

	myToken, err := t.SignedString([]byte(os.Getenv("JWT_SIGNATURE_KEY")))
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:   myToken,
		Expiry:  expiry,
		Created: time.Now(),
	}, nil
}
