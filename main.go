package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"thdr/bstck_2/team"
)

func main() {
	db, err := sqlx.Connect("mysql", "bstack_user:password@(localhost:3306)/bstack_db?tls=false&parseTime=true")
	if err != nil {
		log.Panicf("Error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicf("Failed to ping db: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
  	}))

	log.Printf("Successfully connected to DB")

	tRepo := team.NewRepo(db)
	tHandler := team.NewHandler(tRepo)
	r.Mount("/api/team", tHandler.GetHandler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi there!"))
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Failed to run http server: %v", err)
	}
}
