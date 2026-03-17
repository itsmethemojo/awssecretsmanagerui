package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsmethemojo/awssecretsmanagerui/server/actions"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateJWTAuth() echojwt.Config {
	return echojwt.Config{
		TokenLookup: "header:Authorization:Bearer ",
		SigningKey:  []byte(os.Getenv("JWT_SIGNATURE_KEY")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &actions.LoginClaims{}
		},
		ErrorHandler: JWTErrorChecker,
	}
}

func JWTErrorChecker(c echo.Context, err error) error {
	echoError, ok := err.(*echo.HTTPError)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	return echo.NewHTTPError(http.StatusUnauthorized, echoError.Message)
}
