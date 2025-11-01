package models

import (
	"go_note_pad/config"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetNotes() ([]Note, error) {
	db := config.DBConn()

	rows, err := db.Query("SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func GetNote(id int) (Note, error) {
	db := config.DBConn()

	var note Note
	err := db.QueryRow("SELECT id, title, content, created_at FROM notes WHERE id = ?", id).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func CreateNote(note Note) (Note, error) {
	db := config.DBConn()

	// O campo created_at é definido como padrão pelo banco de dados, então não o incluímos na inserção.
	result, err := db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		return Note{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Note{}, err
	}

	// Busca a nota recém-criada para obter todos os campos, incluindo o ID e CreatedAt gerados.
	return GetNote(int(id))
}

func UpdateNote(note Note) error {
	db := config.DBConn()

	_, err := db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, note.ID)
	return err
}

func DeleteNote(id int) error {
	db := config.DBConn()

	_, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}
