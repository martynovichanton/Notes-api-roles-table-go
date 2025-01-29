package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"notes-api-go/db/database"
	"notes-api-go/middleware"
)

type NoteRoutes struct {
	Queries *database.Queries
}

// Create a new note
func (nr *NoteRoutes) Create(w http.ResponseWriter, r *http.Request) {
	// userID, _ := strconv.Atoi(r.Header.Get("UserID"))
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// var note struct {
	// 	Content string `json:"content"`
	// }
	var note database.CreateNoteParams
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := nr.Queries.CreateNote(context.Background(), database.CreateNoteParams{
		UserID:  userID,
		Content: note.Content,
	}); err != nil {
		http.Error(w, "Could not create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get all notes for the logged-in user
func (nr *NoteRoutes) GetNotesForUser(w http.ResponseWriter, r *http.Request) {
	// userID, _ := strconv.Atoi(r.Header.Get("UserID"))
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notes, err := nr.Queries.GetNotesByUserIDWithUserNames(context.Background(), int64(userID))
	if err != nil {
		http.Error(w, "Could not retrieve notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Get all notes
func (nr *NoteRoutes) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := nr.Queries.GetNotesWithUserNames(context.Background())
	if err != nil {
		http.Error(w, "Could not retrieve notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Update an existing note
func (nr *NoteRoutes) Update(w http.ResponseWriter, r *http.Request) {
	// userID, _ := strconv.Atoi(r.Header.Get("UserID"))
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var updateRequest struct {
		ID      int64  `json:"id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := nr.Queries.UpdateNote(context.Background(), database.UpdateNoteParams{
		Content: updateRequest.Content,
		ID:      updateRequest.ID,
		UserID:  int64(userID),
	}); err != nil {
		http.Error(w, "Could not update note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete an existing note
func (nr *NoteRoutes) Delete(w http.ResponseWriter, r *http.Request) {
	// userID, _ := strconv.Atoi(r.Header.Get("UserID"))
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var deleteRequest struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := nr.Queries.DeleteNote(context.Background(), database.DeleteNoteParams{
		ID:     deleteRequest.ID,
		UserID: int64(userID),
	}); err != nil {
		http.Error(w, "Could not delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
