package repository

import "go_note_pad/models"

type NoteRepository interface {
	Create(note models.Note) (models.Note, error)
	FindAll() ([]models.Note, error)
	FindByID(id int) (models.Note, error)
	Update(note models.Note) error
	Delete(id int) error
}
