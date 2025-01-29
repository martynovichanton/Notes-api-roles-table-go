package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"notes-api-go/db/database"
	"notes-api-go/middleware"
	"notes-api-go/routes"
	"os"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// var DB *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	// db, err := sql.Open("sqlite3", "./db/users.db")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	queries := database.New(db)

	mux := http.NewServeMux()

	authRoutes := &routes.AuthRoutes{Queries: queries}
	userRoutes := &routes.UserRoutes{Queries: queries, DB: db}
	noteRoutes := &routes.NoteRoutes{Queries: queries}

	// routes.DB = DB

	mux.HandleFunc("POST /auth/login", middleware.Logger(authRoutes.Login))
	mux.HandleFunc("POST /auth/refresh", middleware.Logger(authRoutes.RefreshToken))

	// mux.HandleFunc("POST /users", middleware.Logger(userRoutes.CreateUser))
	mux.HandleFunc("POST /users", middleware.Logger(middleware.AuthMiddleware(middleware.AdminMiddleware(userRoutes.CreateUser))))
	mux.HandleFunc("GET /users", middleware.Logger(middleware.AuthMiddleware(middleware.AdminMiddleware(userRoutes.GetUsers))))
	mux.HandleFunc("PUT /users", middleware.Logger(middleware.AuthMiddleware(middleware.AdminMiddleware(userRoutes.UpdateUser))))
	mux.HandleFunc("DELETE /users", middleware.Logger(middleware.AuthMiddleware(middleware.AdminMiddleware(userRoutes.DeleteUser))))

	mux.HandleFunc("POST /notes", middleware.Logger(middleware.AuthMiddleware(noteRoutes.Create)))
	mux.HandleFunc("GET /notes", middleware.Logger(middleware.AuthMiddleware(noteRoutes.GetNotesForUser)))
	mux.HandleFunc("GET /notes/all", middleware.Logger(middleware.AuthMiddleware(middleware.AdminMiddleware(noteRoutes.GetAllNotes))))
	mux.HandleFunc("PUT /notes", middleware.Logger(middleware.AuthMiddleware(noteRoutes.Update)))
	mux.HandleFunc("DELETE /notes", middleware.Logger(middleware.AuthMiddleware(noteRoutes.Delete)))

	log.Println("Server is running on :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
