package models

import (
	"log"
	"time"

	"github.com/natintosh/gowebtutorial/database"
	"github.com/natintosh/gowebtutorial/utils"
)

// Note :
type Note struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Note      string `json:"note"`
	CreatedOn int64  `json:"createdOn,omitempty"`
}

// AddNewNote :
var AddNewNote = func(note Note) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)

	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}
	if note.CreatedOn == 0 {
		note.CreatedOn = utils.GetTimeInMilliseconds(time.Now())
	}

	sqlStatement := `INSERT INTO notes (note_id, user_id, note, created_on) VALUES ($1, $2, $3, $4) RETURNING note_id`

	var id string
	err := db.QueryRow(sqlStatement, note.ID, note.UserID, note.Note, note.CreatedOn).Scan(&id)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to add new note")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := note
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 400
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}

	}

	return result
}

// UpdateNote :
var UpdateNote = func(note Note) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)

	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}
	sqlStatement := `UPDATE notes SET note = $1 WHERE note_id = $2 AND user_id = $3 RETURNING note_id`

	var id string
	err := db.QueryRow(sqlStatement, note.Note, note.ID, note.UserID).Scan(&id)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to update note")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := note
		result = utils.SuccessMessage{Result: "success", Data: data}
	default:
		errobj["code"] = 400
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}

	}

	return result
}

// GetOneNote :
var GetOneNote = func(noteID string) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)
	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}

	sqlStatement := `SELECT * FROM notes WHERE note_id = $1`

	var id string
	var userID string
	var note string
	var createdOn int64

	err := db.QueryRow(sqlStatement, noteID).Scan(&id, &userID, &note, &createdOn)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Note not found")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := Note{ID: id, UserID: userID, Note: note, CreatedOn: createdOn}
		result = utils.SuccessMessage{Result: "success", Data: data}

	default:
		errobj["code"] = 404
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result
}

// GetAllNotes :
var GetAllNotes = func() interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)
	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}

	sqlStatement := `SELECT * FROM notes`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Note not found")
	}
	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		var id string
		var userID string
		var note string
		var createdOn int64

		err = rows.Scan(&id, &userID, &note, &createdOn)
		if err != nil {
			log.Println(err, err.Error())
			errmessage = append(errmessage, "Error occurred while getting data")
		}
		notes = append(notes, Note{ID: id, UserID: userID, Note: note, CreatedOn: createdOn})
	}

	err = rows.Err()
	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Error occurred while getting data")
	}
	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := notes
		result = utils.SuccessMessage{Result: "success", Data: data}

	default:
		errobj["code"] = 500
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result

}

// DeleteNote :
var DeleteNote = func(noteID string) interface{} {
	db := database.GetDb()
	defer db.Close()

	errobj := make(map[string]interface{})
	errmessage := make([]string, 0)

	if err := db.Ping(); err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Failed to connect to database")
	}

	sqlStatement := `DELETE FROM notes WHERE note_id = $1`

	res, err := db.Exec(sqlStatement, noteID)

	if err != nil {
		log.Println(err, err.Error())
		errmessage = append(errmessage, "Error occurred deleting row")
	}

	count, err := res.RowsAffected()

	if count == 0 {
		// log.Println(err, err.Error())
		errmessage = append(errmessage, "No item with id was found or deleted")
	}

	var result interface{}

	switch {
	case len(errmessage) == 0:
		data := make(map[string]interface{})
		data["id"] = noteID
		result = utils.SuccessMessage{Result: "success", Data: data}

	default:
		errobj["code"] = 404
		errobj["message"] = errmessage
		result = utils.ErrorMessage{Result: "error", Error: errobj}
	}

	return result
}
