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

			w.Header().Add("Content-Type", "application/json")
			if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
				w.WriteHeader(statusCode)
			}
			json.NewEncoder(w).Encode(result)
		}
	}
}

// GetAllUsersHandler :
var GetAllUsersHandler = func(w http.ResponseWriter, r *http.Request) {
	result := models.GetAllUsers()

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
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
	if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
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
		if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
			w.WriteHeader(statusCode)
		}
		json.NewEncoder(w).Encode(result)
	default:
		statusCode = 500
		var result interface{}
		errmessage := make([]string, 0)
		errobj := make(map[string]interface{})
		errmessage = append(errmessage, "Unable to delete user", "User ID don't match")
		errobj["code"] = statusCode
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
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
		bytepassword := []byte(user.Password)
		encripted, err := bcrypt.GenerateFromPassword(bytepassword, bcrypt.DefaultCost)
		switch {
		case err != nil:
			log.Println(err, err.Error())
			statusCode = 500
			w.WriteHeader(statusCode)
		default:
			user.Password = string(encripted)
			result := models.UpdateUserPassword(user)

			w.Header().Add("Content-Type", "application/json")
			if exist, statusCode := utils.GetStatusCode(result, utils.ErrorMessage{}); exist && statusCode != 0 {
				w.WriteHeader(statusCode)
			}
			json.NewEncoder(w).Encode(result)
		}
	default:
		statusCode = 500
		var result interface{}
		errmessage := make([]string, 0)
		errobj := make(map[string]interface{})
		errmessage = append(errmessage, "Unable to change password", "User ID don't match")
		errobj["code"] = statusCode
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}
