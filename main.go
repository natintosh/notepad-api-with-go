package main

import (
	"fmt"

	"github.com/natintosh/gowebtutorial/controllers"

	"net/http"

	"github.com/gorilla/mux"
)

func demoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("It worked")
		next.ServeHTTP(w, r)
	})
}

// Routes :
func Routes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.Headers("Content-Type", "application/json")
	subRouter.HandleFunc("/login", controllers.LogUserInHandler).Methods("POST")
	subRouter.HandleFunc("/register", controllers.RegisterUserHandler).Methods("POST")

	userSubrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	userSubrouter.HandleFunc("/users", controllers.GetAllUsersHandler).Methods("GET")
	userSubrouter.HandleFunc("/users/{id}", controllers.GetUserHandler).Methods("GET")
	userSubrouter.HandleFunc("/users/{id}", controllers.DeleteUserHandler).Methods("DELETE")
	userSubrouter.HandleFunc("/users/{id}/password", controllers.UpdateUserPasswordHandler).Methods("PATCH")
	userSubrouter.Use(demoMiddleware, demoMiddleware)

	noteSubrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	noteSubrouter.HandleFunc("/notes", controllers.GetNoteHandler).Methods("GET").Queries("id", "{id}")
	noteSubrouter.HandleFunc("/notes", controllers.PostNoteHandler).Methods("POST")
	noteSubrouter.HandleFunc("/notes", controllers.GetAllNotesHandler).Methods("GET")
	noteSubrouter.HandleFunc("/notes/{id}", controllers.GetNoteHandler).Methods("GET")
	noteSubrouter.HandleFunc("/notes/{id}", controllers.DeleteNoteHandler).Methods("DELETE")
	noteSubrouter.HandleFunc("/notes/{id}", controllers.UpdateNoteHandler).Methods("PATCH")

	return router
}

func main() {
	r := Routes()
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
