package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgshulgin/go_final_project/cmd/config"
	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

func createToken(payload string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": payload, // Subject (user identifier)
		"iss": "todo",  // Issuer
	})
	// token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func (server TaskServer) Authenticate(resp http.ResponseWriter, req *http.Request) {

	auth := dto.Auth{}
	err := json.NewDecoder(req.Body).Decode(&auth)
	if err != nil {
		msg := fmt.Sprintf("Authenticate: ошибка аутентификации, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	var env config.Config
	pwd := env.GetEnvAsString("TODO_PASSWORD", "")
	if len(pwd) == 0 {
		msg := "Authenticate: нет доступа к системе аутентификации"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
	}

	if !strings.EqualFold(auth.Password, pwd) {
		msg := "Authenticate: неверный пароль"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
		return
	}

	h := sha256.New()
	h.Write([]byte(auth.Password))
	s := hex.EncodeToString(h.Sum(nil))
	tokenString, err := createToken(s)

	// типа формируем JWT
	//token := jwt.New(jwt.SigningMethodHS256)
	//signedToken, err := token.SignedString([]byte("my_secret_key"))

	if err != nil {
		//server.log.Errorf("не получилось сформировать JWT")
		msg := "Authenticate: не получилось сформировать JWT"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
		return
	}

	jwt := dto.JWT{Token: tokenString}
	renderJSON(resp, http.StatusOK, jwt)
	server.log.Printf("Authenticate: доступ открыт, token=%s", tokenString)
}

func (server TaskServer) MiddlewareCheckUserAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			var env config.Config
			pwd := env.GetEnvAsString("TODO_PASSWORD", "")
			if len(pwd) == 0 {
				msg := "MiddlewareCheckUserAuth: нет доступа к системе аутентификации"
				server.log.Errorf(msg)
				renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
			}

			//if len(pwd) > 0 {
			var tokenCookie string
			cookie, err := req.Cookie("token")
			if err == nil {
				tokenCookie = cookie.Value
			}
			token, err := verifyToken(tokenCookie)
			if err != nil {
				msg := fmt.Sprintf("MiddlewareCheckUserAuth: Token verification failed: %v", err)
				server.log.Errorf(msg)
				renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
				return
			}

			//}

			h := sha256.New()
			h.Write([]byte(pwd))
			s := hex.EncodeToString(h.Sum(nil))

			s0, err := token.Claims.GetSubject()

			if !strings.EqualFold(s, s0) {
				msg := "MiddlewareCheckUserAuth:неверный пароль"
				server.log.Errorf(msg)
				renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
				return
			}

			server.log.Printf("доступ открыт")

			next.ServeHTTP(resp, req)
		})
	}
}
