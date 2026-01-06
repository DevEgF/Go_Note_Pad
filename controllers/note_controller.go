package controllers

import (
	"encoding/json"
	"go_note_pad/models"
	"go_note_pad/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type NoteController struct {
	Repo repository.NoteRepository
}

func NewNoteController(repo repository.NoteRepository) *NoteController {
	return &NoteController{Repo: repo}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

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

func (c *NoteController) NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.getNotes(w, r)
	case http.MethodPost:
		c.createNote(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func (c *NoteController) NoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID da nota inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getNote(w, r, id)
	case http.MethodPut:
		c.updateNote(w, r, id)
	case http.MethodDelete:
		c.deleteNote(w, r, id)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func (c *NoteController) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := c.Repo.FindAll()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
}

func (c *NoteController) createNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondWithError(w, http.StatusBadRequest, "Carga de solicitação inválida")
		return
	}
	defer r.Body.Close()

	newNote, err := c.Repo.Create(note)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusCreated, newNote)
}

func (c *NoteController) getNote(w http.ResponseWriter, r *http.Request, id int) {
	note, err := c.Repo.FindByID(id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "Nota não encontrada")
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func (c *NoteController) updateNote(w http.ResponseWriter, r *http.Request, id int) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondWithError(w, http.StatusBadRequest, "Carga de solicitação inválida")
		return
	}
	defer r.Body.Close()
	note.ID = id

	if err := c.Repo.Update(note); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func (c *NoteController) deleteNote(w http.ResponseWriter, r *http.Request, id int) {
	if err := c.Repo.Delete(id); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Erro Interno do Servidor")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
