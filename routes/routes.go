package routes

import (
	"database/sql"
	"net/http"

	"notes-api/handlers"
	"notes-api/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {

	r := mux.NewRouter()

	auth := handlers.AuthHandler{DB: db}
	note := handlers.NoteHandler{DB: db}

	// Публичные маршруты (без токена)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	r.HandleFunc("/register", auth.Register).Methods("POST")
	r.HandleFunc("/login", auth.Login).Methods("POST")

	// Защищённые маршруты (нужен JWT токен)
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.Authenticate)

	protected.HandleFunc("/notes", note.CreateNote).Methods("POST")
	protected.HandleFunc("/notes", note.GetNotes).Methods("GET")
	protected.HandleFunc("/notes/{id}", note.GetNote).Methods("GET")
	protected.HandleFunc("/notes/{id}", note.DeleteNote).Methods("DELETE")

	protected.HandleFunc("/notes/{id}", note.UpdateNote).Methods("PUT")

	return r
}
