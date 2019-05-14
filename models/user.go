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

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)

	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to the database")
	}

	sqlStatement := `SELECT * FROM users WHERE user_id = $1`

	var id string
	var username string
	var email string

	err := db.QueryRow(sqlStatement, userID).Scan(&id, &username, &email)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "User not found")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := User{ID: id, Username: username, Email: email}
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 404
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

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
		result = utils.SuccessMessage{Result: "success", Data: data}
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

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)
	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database", err.Error())
	}

	sqlStatement := `UPDATE userspassphrase SET hash_passphrase = $2 WHERE user_id = $1 RETURNING user_id`
	var id string
	err := db.QueryRow(sqlStatement, user.ID, user.Password).Scan(&id)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Unable to change password", err.Error())
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := User{ID: id}
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 500
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result
}

// DeleteUser :
var DeleteUser = func(user User) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)
	if err := db.Ping(); err != nil {
		log.Println(err.Error())
		errmessage = append(errmessage, "Failed to connect to connect to database", err.Error())
	}

	sqlStatement := `DELETE FROM users WHERE user_id = $1`

	res, _ := db.Exec(sqlStatement, user.ID)
	count, _ := res.RowsAffected()

	switch {
	case count != 1:
		errmessage = append(errmessage, "Unable to delete user")
	default:
		sqlStatement = `DELETE FROM userspassphrase WHERE user_id = $1`

		res, _ = db.Exec(sqlStatement, user.ID)
		count, _ := res.RowsAffected()
		if count != 1 {
			errmessage = append(errmessage, "Unable to delete user")
		}
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := User{ID: user.ID}
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 500
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result
}
