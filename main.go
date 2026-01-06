package main

import (
	"database/sql"
	"fmt"
	"go_note_pad/controllers"
	"go_note_pad/repository"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Abrir a conexão com o banco de dados.
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Falha ao abrir conexão com o banco de dados: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// 2. Instanciar o Repositório passando a conexão.
	noteRepo := repository.NewMySQLRepository(db)

	// 3. Instanciar o Controller passando o Repositório.
	noteController := controllers.NewNoteController(noteRepo)

	// 4. Configurar as rotas e subir o servidor.
	http.HandleFunc("/notes", noteController.NotesHandler)
	http.HandleFunc("/notes/", noteController.NoteHandler)

	log.Println("Servidor iniciado em: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
