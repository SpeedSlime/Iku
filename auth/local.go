package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/SpeedSlime/Asahi"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/reply"
	"golang.org/x/crypto/argon2"
)

const (
    tokenSize   	= 32
	tokenExpiration = 24 * 7 * time.Hour
	saltSize   		= 16
    memory      	= 64 * 1024
    iterations  	= 4
    parallelism 	= 4
    keyLength   	= 32
)

func IsTokenExpired(tokenTime time.Time) bool {
    return time.Since(tokenTime) > tokenExpiration
}

func generateToken() (string, error) {
    token := make([]byte, tokenSize)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }
    return base64.RawURLEncoding.Strict().EncodeToString(token), nil
}

func generateSalt() (string, error) {
    salt := make([]byte, saltSize)
    _, err := rand.Read(salt)
    if err != nil {
        return "", err
    }
    return base64.RawURLEncoding.Strict().EncodeToString(salt), nil
}

func hashPassword(password []byte, salt []byte) (string, error) {
    return base64.URLEncoding.Strict().EncodeToString(argon2.IDKey(password, salt, iterations, memory, parallelism, keyLength)), nil
}

func AuthCreatePostRoute(w http.ResponseWriter, r *http.Request) {
	var user User

	r.ParseForm()
	password := r.FormValue("password")
	user.Username = r.FormValue("username")

	if password == "" || user.Username == "" {
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "Password or Username field misssing."), "AuthCreatePostRoute"); return
	}

	if db.Exists(&user) {
		asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "That username is already in use."), "AuthCreatePostRoute"); return  
	}

    salt, err := generateSalt()
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to create a new account, please try again later. 1"), "AuthCreatePostRoute"); return 
    }

    hashedPassword, err := hashPassword([]byte(password), []byte(salt))
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to create a new account, please try again later. 2"), "AuthCreatePostRoute"); return  
    }

	user.Salt = salt
	user.Password = hashedPassword
    err = db.Insert(&user)
    if err != nil {
		fmt.Println(err)
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to create a new account, please try again later. 3"), "AuthCreatePostRoute"); return  
	}

	asahi.Handle(reply.RespondWithResult(w, http.StatusOK, "Successfully registered."), "AuthCreatePostRoute")
}

func AuthLoginPostRoute(w http.ResponseWriter, r *http.Request) {
	var user User

	r.ParseForm()
	password := r.FormValue("password")
	user.Username = r.FormValue("username")

	if password == "" || user.Username == "" {
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "Password or Username field misssing."), "AuthLoginPostRoute"); return
	}

    found, err := db.Select(&user)
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to login, please try again later."), "AuthLoginPostRoute"); return 
    }

    if !found {
        asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "That account does not exist."), "AuthLoginPostRoute"); return  
    }

    hashedPassword, err := hashPassword([]byte(password), []byte(user.Salt))
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to login, please try again later."), "AuthLoginPostRoute"); return  
    }

    if user.Password != hashedPassword {
        asahi.Handle(reply.RespondWithResult(w, http.StatusUnauthorized, "Invalid password."), "AuthLoginPostRoute"); return 
    }

    token, err := generateToken()
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to login, please try again later."), "AuthLoginPostRoute"); return   
    }

	cond := user
    user.Token = token
	user.TokenTime = time.Now().Add(tokenExpiration)
    _, err = db.Update(&user, &cond)
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to login, please try again later."), "AuthLoginPostRoute"); return  
    }

    asahi.Handle(reply.RespondWithJSON(w, http.StatusOK, TokenRequest{Token: token, Username: user.Username}), "AuthLoginPostRoute")
}

func AuthLogoutPostRoute(w http.ResponseWriter, r *http.Request) {
	var user User

	r.ParseForm()
	user.Token = r.FormValue("token")
	user.Username = r.FormValue("username")

	if user.Token == "" || user.Username == "" {
		asahi.Handle(reply.RespondWithResult(w, http.StatusForbidden, "Token or Username field misssing."), "AuthLogoutPostRoute"); return
	}

    found, err := db.Select(&user)
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to logout, please try again later."), "AuthLogoutPostRoute"); return
    }

    if !found {
        asahi.Handle(reply.RespondWithResult(w, http.StatusBadRequest, "Invalid request."), "AuthLogoutPostRoute"); return
    }

	token, err := generateToken()
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to logout, please try again later."), "AuthLogoutPostRoute"); return 
    }

	cond := user
	user.Token = token
	user.TokenTime = time.Now().Add(-tokenExpiration)
    _, err = db.Update(&user, &cond)
    if err != nil {
        asahi.Handle(reply.RespondWithResult(w, http.StatusInternalServerError, "An unexpected error has occured when trying to logout, please try again later."), "AuthLogoutPostRoute"); return 
    }
    
	asahi.Handle(reply.RespondWithResult(w, http.StatusOK, "Successfully logged out."), "AuthLoginPostRoute")
}