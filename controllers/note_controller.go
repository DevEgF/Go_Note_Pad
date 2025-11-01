package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"go_note_pad/models"
)

var tmpl = template.Must(template.ParseGlob("views/*.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	notes, err := models.GetNotes()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", notes)
}

func NewNote(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "new.html", nil)
}

func SaveNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	note := models.Note{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	err := models.CreateNote(note)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := models.GetNote(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "edit.html", note)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/notes/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note := models.Note{
		ID: 	 id,
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	err = models.UpdateNote(note)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteNote(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
