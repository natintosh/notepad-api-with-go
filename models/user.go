package models

import (
	"log"

	"github.com/natintosh/gowebtutorial/database"
	"github.com/natintosh/gowebtutorial/utils"
)

// User :
type User struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// AddNewUser :
var AddNewUser = func(user User) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)

	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}

	sqlStatement := `INSERT INTO users (user_id, username, email) VALUES ($1, $2, $3) RETURNING user_id`

	var id string
	err := db.QueryRow(sqlStatement, user.ID, user.Username, user.Email).Scan(&id)
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to add new user")
	}

	sqlStatement = `INSERT INTO userspassphrase (user_id, hash_passphrase) VALUES ($1, $2) RETURNING user_id`

	err = db.QueryRow(sqlStatement, user.ID, user.Password).Scan(&id)
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to add new user")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := User{ID: user.ID, Username: user.Username, Email: user.Email}
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 400
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}

	}

	return result
}

// GetOneUser :
var GetOneUser = func(userID string) interface{} {
	db := database.GetDb()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	sqlStatement := `SELECT * FROM users WHERE user_id = $1`

	var id string
	var username string
	var email string

	err := db.QueryRow(sqlStatement, userID).Scan(&id, &username, &email)

	var result map[string]interface{} = make(map[string]interface{})
	var user map[string]interface{} = make(map[string]interface{})

	user["id"] = id
	user["username"] = username
	user["email"] = email

	switch {
	case err == nil:
		result["result"] = "success"
	default:
		result["result"] = "error"
	}

	result["user"] = user

	return result
}

// GetAllUsers :
var GetAllUsers = func() interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)
	if err := db.Ping(); err != nil {
		log.Println(err.Error())
		errmessage = append(errmessage, "Failed to connect to database", err.Error())
	}

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Users not found", err.Error())
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var id string
		var username string
		var email string

		err = rows.Scan(&id, &username, &email)
		if err != nil {
			log.Println(err, err.Error())
			errmessage = append(errmessage, "Error occurred while getting data", err.Error())
		}
		users = append(users, User{ID: id, Username: username, Email: email})
	}

	err = rows.Err()
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Error occurred while getting data", err.Error())
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := users
		result = utils.SuccessMessage{Result: "Success", Data: data}
	default:
		errobj["code"] = 500
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result
}

// UpdateUserPassword :
var UpdateUserPassword = func(user User) interface{} {
	db := database.GetDb()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	sqlStatement := `UPDATE userspassphrase SET hash_passphrase = $2 WHERE user_id = $1 RETURNING user_id`
	var id string
	err := db.QueryRow(sqlStatement, user.ID, user.Password).Scan(&id)

	var result map[string]interface{} = make(map[string]interface{})

	if err != nil {
		panic(err)
	}

	switch {
	case err == nil:
		result["result"] = "success"
	default:
		result["result"] = "error"
	}
	result["id"] = id

	return result
}

// DeleteUser :
var DeleteUser = func(user User) interface{} {
	db := database.GetDb()
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	sqlStatement := `DELETE FROM users WHERE user_id = $1`

	_, err := db.Exec(sqlStatement, user.ID)

	var result map[string]interface{} = make(map[string]interface{})
	if err != nil {
		panic(err)
	}
	switch {
	case err == nil:
		result["result"] = "success"
	default:
		result["result"] = "error"
	}
	sqlStatement = `DELETE FROM userspassphrase WHERE user_id = $1`

	_, err = db.Exec(sqlStatement, user.ID)
	if err != nil {
		panic(err)
	}
	switch {
	case err == nil:
		result["result"] = "success"
	default:
		result["result"] = "error"
	}

	result["id"] = user.ID
	return result
}