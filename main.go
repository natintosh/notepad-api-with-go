package main

import (
	"github.com/natintosh/gowebtutorial/controllers"

	"net/http"

	"github.com/gorilla/mux"
)

// Routes :
func Routes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.Headers("Content-Type", "application/json")
	subrouter.HandleFunc("/users", controllers.GetUserHandler).Methods("GET").Queries("id", "{id}")
	subrouter.HandleFunc("/users", controllers.PostUserHandler).Methods("POST")
	subrouter.HandleFunc("/users", controllers.GetAllUsersHandler).Methods("GET")
	subrouter.HandleFunc("/users/{id}", controllers.GetUserHandler).Methods("GET")
	subrouter.HandleFunc("/users/{id}", controllers.DeleteUserHandler).Methods("DELETE")
	subrouter.HandleFunc("/users/{id}/password", controllers.UpdateUserPasswordHandler).Methods("PATCH")

	subrouter.HandleFunc("/notes", controllers.GetNoteHandler).Methods("GET").Queries("id", "{id}")
	subrouter.HandleFunc("/notes", controllers.PostNoteHandler).Methods("POST")
	subrouter.HandleFunc("/notes", controllers.GetAllNotesHandler).Methods("GET")
	subrouter.HandleFunc("/notes/{id}", controllers.GetNoteHandler).Methods("GET")
	subrouter.HandleFunc("/notes/{id}", controllers.DeleteNoteHandler).Methods("DELETE")
	subrouter.HandleFunc("/notes/{id}", controllers.UpdateNoteHandler).Methods("PATCH")

	return router
}

func main() {
	r := Routes()
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
