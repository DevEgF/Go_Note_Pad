package main

import (
	"log"
	"net/http"
	"go_note_pad/config"
	"go_note_pad/controllers"
)

func main() {
	err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/notes/new", controllers.NewNote)
	http.HandleFunc("/notes/save", controllers.SaveNote)
	http.HandleFunc("/notes/edit/", controllers.EditNote)
	http.HandleFunc("/notes/update/", controllers.UpdateNote)
	http.HandleFunc("/notes/delete/", controllers.DeleteNote)

	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
