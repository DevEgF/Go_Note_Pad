package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go_note_pad/config"
	"go_note_pad/models"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupDB inicializa um mock de banco de dados e o injeta no pacote de configuração.
func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	config.SetDB(db)
	return db, mock
}

func TestGetNotes_Controller(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(1, "Nota 1", "Conteúdo 1", time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")).
		WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
	rr := httptest.NewRecorder()

	NotesHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNote_Controller(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	note := models.Note{Title: "Nova Nota", Content: "Novo Conteúdo"}
	noteJSON, _ := json.Marshal(note)
	newID := int64(1)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO notes (title, content) VALUES (?, ?)")).
		WithArgs(note.Title, note.Content).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(newID, note.Title, note.Content, time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes WHERE id = ?")).
		WithArgs(newID).
		WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(noteJSON))
	rr := httptest.NewRecorder()

	NotesHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdNote models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &createdNote)
	require.NoError(t, err)
	assert.Equal(t, int(newID), createdNote.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNote_Controller(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteID := 1
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(noteID, "Nota de Teste", "Conteúdo de Teste", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes WHERE id = ?")).
		WithArgs(noteID).
		WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
	rr := httptest.NewRecorder()

	NoteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateNote_Controller(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteID := 1
	note := models.Note{Title: "Título Atualizado", Content: "Conteúdo Atualizado"}
	noteJSON, _ := json.Marshal(note)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE notes SET title = ?, content = ? WHERE id = ?")).
		WithArgs(note.Title, note.Content, noteID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ := http.NewRequest(http.MethodPut, "/notes/1", bytes.NewBuffer(noteJSON))
	rr := httptest.NewRecorder()

	NoteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteNote_Controller(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteID := 1
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM notes WHERE id = ?")).
		WithArgs(noteID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ := http.NewRequest(http.MethodDelete, "/notes/1", nil)
	rr := httptest.NewRecorder()

	NoteHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
