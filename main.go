package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"

	"thdr/bstck_2/team"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := sqlx.Connect("mysql", "bstack_user:password@(localhost:3306)/bstack_db?tls=false")
	if err != nil {
		log.Panicf("Error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicf("Failed to ping db: %v", err)
	}

	log.Printf("Successfully connected to DB")

	tRepo := team.NewRepo(db)
	tHandler := team.NewHandler(tRepo)
	r.Mount("/api/team", tHandler.GetHandler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi there!"))
	})

	http.ListenAndServe(":8080", r)
}
