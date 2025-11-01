package models

import (
	"database/sql"
	"go_note_pad/config"
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
	config.SetDB(db) // Injeta o mock no pacote config
	return db, mock
}

func TestGetNotes(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	// Define as linhas que o mock deve retornar
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(1, "Nota 1", "Conteúdo 1", time.Now()).
		AddRow(2, "Nota 2", "Conteúdo 2", time.Now())

	// Define a expectativa para a consulta
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")).
		WillReturnRows(rows)

	// Chama a função a ser testada
	notes, err := GetNotes()

	// Verifica os resultados
	require.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, "Nota 1", notes[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNote(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteID := 1
	expectedNote := Note{
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

	note, err := GetNote(noteID)

	require.NoError(t, err)
	assert.Equal(t, expectedNote, note)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNote(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteToCreate := Note{Title: "Nova Nota", Content: "Novo Conteúdo"}
	newID := int64(1)

	// Expectativa para a inserção
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO notes (title, content) VALUES (?, ?)")).
		WithArgs(noteToCreate.Title, noteToCreate.Content).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	// Expectativa para a busca da nota recém-criada
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(newID, noteToCreate.Title, noteToCreate.Content, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at FROM notes WHERE id = ?")).
		WithArgs(newID).
		WillReturnRows(rows)

	createdNote, err := CreateNote(noteToCreate)

	require.NoError(t, err)
	assert.Equal(t, int(newID), createdNote.ID)
	assert.Equal(t, noteToCreate.Title, createdNote.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateNote(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteToUpdate := Note{ID: 1, Title: "Título Atualizado", Content: "Conteúdo Atualizado"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE notes SET title = ?, content = ? WHERE id = ?")).
		WithArgs(noteToUpdate.Title, noteToUpdate.Content, noteToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada

	err := UpdateNote(noteToUpdate)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteNote(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	noteID := 1

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM notes WHERE id = ?")).
		WithArgs(noteID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada

	err := DeleteNote(noteID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
