package repository

import (
	"database/sql"
	"go_note_pad/models"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) FindAll() ([]models.Note, error) {
	rows, err := r.db.Query("SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *MySQLRepository) FindByID(id int) (models.Note, error) {
	var note models.Note
	err := r.db.QueryRow("SELECT id, title, content, created_at FROM notes WHERE id = ?", id).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (r *MySQLRepository) Create(note models.Note) (models.Note, error) {
	result, err := r.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		return models.Note{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Note{}, err
	}

	return r.FindByID(int(id))
}

func (r *MySQLRepository) Update(note models.Note) error {
	_, err := r.db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, note.ID)
	return err
}

func (r *MySQLRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}
