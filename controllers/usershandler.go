package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/natintosh/gowebtutorial/models"
	"github.com/natintosh/gowebtutorial/utils"
	"github.com/rs/xid"

	"github.com/gorilla/mux"
)

// PostUserHandler :
var PostUserHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	statusCode := 0
	switch {
	case err != nil:
		log.Println(err, err.Error())
		statusCode = 500
		w.WriteHeader(statusCode)
	default:
		gxid := xid.New().String()
		user.ID = gxid
		bytepassword := []byte(user.Password)
		encripted, err := bcrypt.GenerateFromPassword(bytepassword, bcrypt.DefaultCost)
		switch {
		case err != nil:
			log.Println(err, err.Error())
			statusCode = 500
			w.WriteHeader(statusCode)
		default:
			user.Password = string(encripted)
			result := models.AddNewUser(user)
			if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
				w.WriteHeader(statusCode)
			}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	}
}

// GetAllUsersHandler :
var GetAllUsersHandler = func(w http.ResponseWriter, r *http.Request) {
	result := models.GetAllUsers()

	if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetUserHandler :
var GetUserHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := models.GetOneUser(id)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeleteUserHandler :
var DeleteUserHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	var getUser models.User
	err := decoder.Decode(&getUser)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{} = models.DeleteUser(getUser).(map[string]interface{})

	if id != result["id"] {
		result["result"] = "error"
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

// UpdateUserPasswordHandler :
var UpdateUserPasswordHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	var getUser models.User
	err := decoder.Decode(&getUser)
	if err != nil {
		panic(err)
	}
	bytepassword := []byte(getUser.Password)
	encripted, err := bcrypt.GenerateFromPassword(bytepassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	getUser.Password = string(encripted)

	var result map[string]interface{} = models.UpdateUserPassword(getUser).(map[string]interface{})

	if id != result["id"] {
		result["result"] = "error"
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
