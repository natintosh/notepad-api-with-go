package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"crypto/sha256"

	"github.com/natintosh/gowebtutorial/models"
	"github.com/natintosh/gowebtutorial/utils"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

// LogUserInHandler :
var LogUserInHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	switch {
	case err != nil:
		log.Println(err, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		userDoesExist := models.LogUserIn(user)

		if userDoesExist {

			type Claims struct {
				Username string `json:"username"`
				jwt.StandardClaims
			}

			secretKey := os.Getenv("secret_key")
			h := sha256.New()
			h.Write([]byte(secretKey))
			var jwtKey = []byte(h.Sum(nil))
			expirationTime := time.Now().Add(3600 * time.Second)

			claims := &Claims{
				Username: user.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				log.Println(err, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			tokenResult := map[string]interface{}{"token": tokenString}
			json.NewEncoder(w).Encode(tokenResult)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}
}

// RegisterUserHandler :
var RegisterUserHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	switch {
	case err != nil:
		log.Println(err, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		gxid := xid.New().String()
		user.ID = gxid
		bytePassword := []byte(user.Password)
		encripted, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
		switch {
		case err != nil:
			log.Println(err, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			user.Password = string(encripted)
			payload := models.AddNewUser(user)

			if exist, statusCode := utils.GetStatusCode(payload); exist && statusCode != 0 {
				w.WriteHeader(statusCode)
				w.Header().Add("Content-Type", "application/json")
				json.NewEncoder(w).Encode(payload)
				return
			}

			type Claims struct {
				Username string `json:"username"`
				jwt.StandardClaims
			}

			secretKey := os.Getenv("secret_key")
			h := sha256.New()
			h.Write([]byte(secretKey))
			var jwtKey = []byte(h.Sum(nil))
			expirationTime := time.Now().Add(3600 * time.Second)

			claims := &Claims{
				Username: user.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)

			if err != nil {
				log.Println(err, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tokenResult := map[string]interface{}{"token": tokenString, "data": payload}

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResult)
		}
	}
}

// GetAllUsersHandler :
var GetAllUsersHandler = func(w http.ResponseWriter, r *http.Request) {
	result := models.GetAllUsers()

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	json.NewEncoder(w).Encode(result)
}

// GetUserHandler :
var GetUserHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := models.GetOneUser(id)

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	json.NewEncoder(w).Encode(result)
}

// DeleteUserHandler :
var DeleteUserHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	statusCode := 0
	switch {
	case err != nil:
		log.Println(err, err.Error())
		statusCode = 500
		w.WriteHeader(statusCode)
	case id == user.ID:
		result := models.DeleteUser(user)

		w.Header().Add("Content-Type", "application/json")
		if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
			w.WriteHeader(statusCode)
		}
		json.NewEncoder(w).Encode(result)
	default:
		statusCode = 500
		var result interface{}
		errMessage := make([]string, 0)
		errObj := make(map[string]interface{})
		errMessage = append(errMessage, "Unable to delete user", "User ID don't match")
		errObj["code"] = statusCode
		errObj["message"] = errMessage
		result = utils.ErrorMessage{Result: "error", Error: errObj}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}

// UpdateUserPasswordHandler :
var UpdateUserPasswordHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	statusCode := 0
	switch {
	case err != nil:
		log.Println(err, err.Error())
		statusCode = 500
		w.WriteHeader(statusCode)
	case id == user.ID:
		bytePassword := []byte(user.Password)
		encripted, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
		switch {
		case err != nil:
			log.Println(err, err.Error())
			statusCode = 500
			w.WriteHeader(statusCode)
		default:
			user.Password = string(encripted)
			result := models.UpdateUserPassword(user)

			w.Header().Add("Content-Type", "application/json")
			if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
				w.WriteHeader(statusCode)
			}
			json.NewEncoder(w).Encode(result)
		}
	default:
		statusCode = 500
		var result interface{}
		errMessage := make([]string, 0)
		errObj := make(map[string]interface{})
		errMessage = append(errMessage, "Unable to change password", "User ID don't match")
		errObj["code"] = statusCode
		errObj["message"] = errMessage
		result = utils.ErrorMessage{Result: "error", Error: errObj}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}
