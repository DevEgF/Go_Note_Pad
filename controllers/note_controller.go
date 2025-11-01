package controllers

import (
	"encoding/json"
	"go_note_pad/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// respondWithError é uma função auxiliar para escrever uma resposta de erro JSON.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON é uma função auxiliar para escrever uma resposta JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falha ao serializar JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// NotesHandler lida com todas as solicitações para /notes
func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNotes(w, r)
	case http.MethodPost:
		createNote(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

// NoteHandler lida com todas as solicitações para /notes/{id}
func NoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID da nota inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		getNote(w, r, id)
	case http.MethodPut:
		updateNote(w, r, id)
	case http.MethodDelete:
		deleteNote(w, r, id)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := models.GetNotes()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondWithError(w, http.StatusBadRequest, "Carga de solicitação inválida")
		return
	}
	defer r.Body.Close()

	newNote, err := models.CreateNote(note)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusCreated, newNote)
}

func getNote(w http.ResponseWriter, r *http.Request, id int) {
	note, err := models.GetNote(id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "Nota não encontrada")
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func updateNote(w http.ResponseWriter, r *http.Request, id int) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondWithError(w, http.StatusBadRequest, "Carga de solicitação inválida")
		return
	}
	defer r.Body.Close()
	note.ID = id

	if err := models.UpdateNote(note); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func deleteNote(w http.ResponseWriter, r *http.Request, id int) {
	if err := models.DeleteNote(id); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
