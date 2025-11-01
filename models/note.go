package models

import (
	"time"
	"go_note_pad/config"
)

type Note struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
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
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func GetNote(id int) (Note, error) {
	db := config.DBConn()

	var note Note
	err := db.QueryRow("SELECT id, title, content, created_at FROM notes WHERE id = ?", id).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func CreateNote(note Note) error {
	db := config.DBConn()

	_, err := db.Exec("INSERT INTO notes (title, content, created_at) VALUES (?, ?, ?)", note.Title, note.Content, time.Now())
	return err
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
