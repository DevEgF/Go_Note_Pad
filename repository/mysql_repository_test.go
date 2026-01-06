package repository

import (
	"go_note_pad/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) (*MySQLRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := NewMySQLRepository(db)
	return repo, mock
}

func TestFindAll(t *testing.T) {
	repo, mock := setupRepo(t)
	defer repo.db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(1, "Nota 1", "Conteúdo 1", time.Now()).
		AddRow(2, "Nota 2", "Conteúdo 2", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")).
		WillReturnRows(rows)

	notes, err := repo.FindAll()

	require.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, "Nota 1", notes[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID(t *testing.T) {
	repo, mock := setupRepo(t)
	defer repo.db.Close()

	noteID := 1
	expectedNote := models.Note{
		ID:        noteID,
		Title:     "Nota de Teste",
		Content:   "Conteúdo de Teste",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(expectedNote.ID, expectedNote.Title, expectedNote.Content, expectedNote.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes WHERE id = ?")).
		WithArgs(noteID).
		WillReturnRows(rows)

	note, err := repo.FindByID(noteID)

	require.NoError(t, err)
	assert.Equal(t, expectedNote, note)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	repo, mock := setupRepo(t)
	defer repo.db.Close()

	noteToCreate := models.Note{Title: "Nova Nota", Content: "Novo Conteúdo"}
	newID := int64(1)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO notes (title, content) VALUES (?, ?)")).
		WithArgs(noteToCreate.Title, noteToCreate.Content).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(newID, noteToCreate.Title, noteToCreate.Content, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes WHERE id = ?")).
		WithArgs(newID).
		WillReturnRows(rows)

	createdNote, err := repo.Create(noteToCreate)

	require.NoError(t, err)
	assert.Equal(t, int(newID), createdNote.ID)
	assert.Equal(t, noteToCreate.Title, createdNote.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	repo, mock := setupRepo(t)
	defer repo.db.Close()

	noteToUpdate := models.Note{ID: 1, Title: "Título Atualizado", Content: "Conteúdo Atualizado"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE notes SET title = ?, content = ? WHERE id = ?")).
		WithArgs(noteToUpdate.Title, noteToUpdate.Content, noteToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(noteToUpdate)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	repo, mock := setupRepo(t)
	defer repo.db.Close()

	noteID := 1

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM notes WHERE id = ?")).
		WithArgs(noteID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(noteID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
