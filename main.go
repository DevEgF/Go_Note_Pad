package main

import (
	"go_note_pad/config"
	"go_note_pad/controllers"
	"log"
	"net/http"
)

func main() {
	// Inicializa o banco de dados e lida com erros de conexão
	if err := config.InitDB(); err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Define os manipuladores para os endpoints da API
	http.HandleFunc("/notes", controllers.NotesHandler)
	http.HandleFunc("/notes/", controllers.NoteHandler)

	// Inicia o servidor na porta 8080
	log.Println("Servidor iniciado em: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
