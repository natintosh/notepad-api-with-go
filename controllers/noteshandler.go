package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/natintosh/gowebtutorial/models"
	"github.com/natintosh/gowebtutorial/utils"
)

// GetNoteHandler :
var GetNoteHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := models.GetOneNote(id)

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	json.NewEncoder(w).Encode(result)

}

// GetAllNotesHandler :
var GetAllNotesHandler = func(w http.ResponseWriter, r *http.Request) {
	result := models.GetAllNotes()

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	json.NewEncoder(w).Encode(result)
}

// PostNoteHandler :
var PostNoteHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var note models.Note
	err := decoder.Decode(&note)
	statusCode := 0

	switch {
	case err != nil:
		log.Println(err, err.Error())
		statusCode = 500
		w.WriteHeader(statusCode)
	default:
		result := models.AddNewNote(note)
		w.Header().Add("Content-Type", "application/json")
		if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
			w.WriteHeader(statusCode)
		}
		json.NewEncoder(w).Encode(result)

	}
}

// UpdateNoteHandler :
var UpdateNoteHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var note models.Note
	err := decoder.Decode(&note)
	statusCode := 0

	switch {
	case err != nil:
		log.Println(err, err.Error())
		statusCode = 400
		w.WriteHeader(statusCode)
	default:
		result := models.UpdateNote(note)
		w.Header().Add("Content-Type", "application/json")
		if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
			w.WriteHeader(statusCode)
		}
		json.NewEncoder(w).Encode(result)

	}
}

// DeleteNoteHandler :
var DeleteNoteHandler = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := models.DeleteNote(id)

	w.Header().Add("Content-Type", "application/json")
	if exist, statusCode := utils.GetStatusCode(result); exist && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	json.NewEncoder(w).Encode(result)
}
