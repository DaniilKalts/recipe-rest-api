package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DaniilKalts/recipe-rest-api/api/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Create the connection string
	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"

	// Connect to the Database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handlers.InitializeDB(db)

	r := mux.NewRouter()

	// GET endpoint to get recipes
	r.HandleFunc("/recipes", handlers.GetRecipes).Methods("GET")

	// DELETE endpoint to remove a recipe
	r.HandleFunc("/recipes", handlers.DeleteRecipe).Methods("DELETE")

	// POST endpoint to create a new recipe
	r.HandleFunc("/recipes", handlers.CreateRecipe).Methods("POST")

	// UPDATE endpoint to update an existing recipe
	r.HandleFunc("/recipes", handlers.UpdateRecipe).Methods("PUT")

	// Start the local server
	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}