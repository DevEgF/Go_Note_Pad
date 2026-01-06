package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_note_pad/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockNoteRepository is a mock implementation of repository.NoteRepository
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Create(note models.Note) (models.Note, error) {
	args := m.Called(note)
	return args.Get(0).(models.Note), args.Error(1)
}

func (m *MockNoteRepository) FindAll() ([]models.Note, error) {
	args := m.Called()
	return args.Get(0).([]models.Note), args.Error(1)
}

func (m *MockNoteRepository) FindByID(id int) (models.Note, error) {
	args := m.Called(id)
	return args.Get(0).(models.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(note models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

func (m *MockNoteRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupController(t *testing.T) (*NoteController, *MockNoteRepository) {
	mockRepo := new(MockNoteRepository)
	controller := NewNoteController(mockRepo)
	return controller, mockRepo
}

func TestGetNotes_Controller(t *testing.T) {
	controller, mockRepo := setupController(t)

	notes := []models.Note{
		{ID: 1, Title: "Nota 1", Content: "Conteúdo 1", CreatedAt: time.Now()},
		{ID: 2, Title: "Nota 2", Content: "Conteúdo 2", CreatedAt: time.Now()},
	}

	mockRepo.On("FindAll").Return(notes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
	rr := httptest.NewRecorder()

	controller.NotesHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetNotes_Controller_Error(t *testing.T) {
	controller, mockRepo := setupController(t)

	mockRepo.On("FindAll").Return([]models.Note{}, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
	rr := httptest.NewRecorder()

	controller.NotesHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateNote_Controller(t *testing.T) {
	controller, mockRepo := setupController(t)

	noteInput := models.Note{Title: "Nova Nota", Content: "Novo Conteúdo"}
	createdNote := models.Note{ID: 1, Title: "Nova Nota", Content: "Novo Conteúdo", CreatedAt: time.Now()}

	// Note: createdNote will differ from noteInput (ID=0 vs ID=1), so match based on fields if strictly necessary,
	// but here we just pass the object coming from json decode.
	// Since json decode creates a clean struct, ID is 0.
	mockRepo.On("Create", noteInput).Return(createdNote, nil)

	noteJSON, _ := json.Marshal(noteInput)
	req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(noteJSON))
	rr := httptest.NewRecorder()

	controller.NotesHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var respNote models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &respNote)
	require.NoError(t, err)
	assert.Equal(t, createdNote.ID, respNote.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetNote_Controller(t *testing.T) {
	controller, mockRepo := setupController(t)

	noteID := 1
	note := models.Note{ID: noteID, Title: "Nota de Teste", Content: "Conteúdo de Teste", CreatedAt: time.Now()}

	mockRepo.On("FindByID", noteID).Return(note, nil)

	req, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
	rr := httptest.NewRecorder()

	controller.NoteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateNote_Controller(t *testing.T) {
	controller, mockRepo := setupController(t)

	noteID := 1
	noteInput := models.Note{Title: "Título Atualizado", Content: "Conteúdo Atualizado"}
	// controller sets ID manually before calling update
	expectedNoteCall := noteInput
	expectedNoteCall.ID = noteID

	mockRepo.On("Update", expectedNoteCall).Return(nil)

	noteJSON, _ := json.Marshal(noteInput)
	req, _ := http.NewRequest(http.MethodPut, "/notes/1", bytes.NewBuffer(noteJSON))
	rr := httptest.NewRecorder()

	controller.NoteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteNote_Controller(t *testing.T) {
	controller, mockRepo := setupController(t)

	noteID := 1
	mockRepo.On("Delete", noteID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/notes/1", nil)
	rr := httptest.NewRecorder()

	controller.NoteHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}
