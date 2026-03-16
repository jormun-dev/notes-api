package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteHandler struct {
	DB *sql.DB
}

// Вспомогательная функция — получить username из контекста
func getUsernameFromContext(r *http.Request) string {
	username, _ := r.Context().Value("username").(string)
	return username
}

func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {

	username := getUsernameFromContext(r)

	rows, err := h.DB.Query(`
		SELECT n.id, n.title, n.content 
		FROM notes n
		JOIN users u ON n.user_id = u.id
		WHERE u.username = $1
	`, username)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		rows.Scan(&n.ID, &n.Title, &n.Content)
		notes = append(notes, n)
	}

	// Вернуть пустой массив вместо null
	if notes == nil {
		notes = []Note{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	username := getUsernameFromContext(r)

	var title, content string

	err := h.DB.QueryRow(`
		SELECT n.title, n.content 
		FROM notes n
		JOIN users u ON n.user_id = u.id
		WHERE n.id = $1 AND u.username = $2
	`, id, username).Scan(&title, &content)

	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"title":   title,
		"content": content,
	})
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {

	username := getUsernameFromContext(r)

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var id int
	err := h.DB.QueryRow(`
		INSERT INTO notes(title, content, user_id)
		VALUES($1, $2, (SELECT id FROM users WHERE username=$3))
		RETURNING id
	`, note.Title, note.Content, username).Scan(&id)

	if err != nil {
		http.Error(w, "Could not create note", http.StatusInternalServerError)
		return
	}

	note.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	username := getUsernameFromContext(r)

	result, err := h.DB.Exec(`
		DELETE FROM notes 
		WHERE id = $1 
		AND user_id = (SELECT id FROM users WHERE username = $2)
	`, id, username)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	username := getUsernameFromContext(r)

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE notes SET title = $1, content = $2
		WHERE id = $3
		AND user_id = (SELECT id FROM users WHERE username = $4)
	`, note.Title, note.Content, id, username)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note updated",
	})
}
