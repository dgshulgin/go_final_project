package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dgshulgin/go_final_project/cmd/config"
	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

func createToken(payload string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": payload,
		"iss": "todo",
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidJWToken
	}

	return token, nil
}

func (server TaskServer) Authenticate(resp http.ResponseWriter, req *http.Request) {

	auth := dto.Auth{}
	err := json.NewDecoder(req.Body).Decode(&auth)
	if err != nil {
		msg := errors.Join(ErrAuthentication, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	var env config.Config
	pwd := env.GetEnvAsString("TODO_PASSWORD", "")
	if len(pwd) == 0 {
		msg := ErrAuthentication.Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
	}

	if !strings.EqualFold(auth.Password, pwd) {
		msg := ErrWrongPassword.Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
		return
	}

	h := sha256.New()
	h.Write([]byte(auth.Password))
	s := hex.EncodeToString(h.Sum(nil))
	tokenString, err := createToken(s)
	if err != nil {
		msg := errors.Join(ErrCreateJWT, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
		return
	}

	// Успех, доступ открыт
	jwt := dto.JWT{Token: tokenString}
	renderJSON(resp, http.StatusOK, jwt)
}

func (server TaskServer) MiddlewareCheckUserAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			var env config.Config
			pwd := env.GetEnvAsString("TODO_PASSWORD", "")
			if len(pwd) == 0 {
				msg := ErrAuthentication.Error()
				server.logging(msg, nil)
				renderJSON(resp,
					http.StatusInternalServerError, dto.Error{Error: msg})
			}

			var tokenCookie string
			cookie, err := req.Cookie("token")
			if err == nil {
				tokenCookie = cookie.Value
			}
			token, err := verifyToken(tokenCookie)
			if err != nil {
				msg := errors.Join(ErrTokenVerification, err).Error()
				server.logging(msg, nil)
				renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
				return
			}

			h := sha256.New()
			h.Write([]byte(pwd))
			s := hex.EncodeToString(h.Sum(nil))

			s0, err := token.Claims.GetSubject()
			if !strings.EqualFold(s, s0) {
				msg := ErrWrongPassword.Error()
				server.logging(msg, nil)
				renderJSON(resp, http.StatusUnauthorized, dto.Error{Error: msg})
				return
			}

			next.ServeHTTP(resp, req)
		})
	}
}
