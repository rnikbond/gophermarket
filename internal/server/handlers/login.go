package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var users map[string]string

func (a Auth) String() string {
	return fmt.Sprintf("\tlogin: %s\n\tpassword: %s\n", a.Login, a.Password)
}

// Login - Авторизация пользователя
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			// @TODO :: Может http.StatusUnsupportedMediaType ?
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("error close body: %v\n", err)
			}
		}()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error read body: %v\n", err)
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}

		var auth Auth
		if err := json.Unmarshal(data, &auth); err != nil {
			log.Printf("error unmarshal json: %v\n", err)
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		pass, ok := users[auth.Login]
		if !ok || pass != auth.Password {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Register - Регистрация пользователя
func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			// @TODO :: Может http.StatusUnsupportedMediaType ?
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("error close body: %v\n", err)
			}
		}()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error read body: %v\n", err)
			// @TODO :: Может тут что-то надо написать в ответ?
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var auth Auth
		if err := json.Unmarshal(data, &auth); err != nil {
			log.Printf("error unmarshal JSON: %v\n", err)
			http.Error(w, "invalid JSON", http.StatusBadRequest)
		}

		if users == nil {
			users = make(map[string]string)
		}

		if _, ok := users[auth.Login]; ok {
			http.Error(w, "already exists", http.StatusConflict)
			return
		}

		users[auth.Login] = auth.Password
	}
}
