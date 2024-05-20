package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Recipe struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Ingredients  string    `json:"ingredients"`
	Instructions string    `json:"instructions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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

	r := mux.NewRouter()

	// GET endpoint to get recipes
	r.HandleFunc("/recipes", getRecipes).Methods("GET")

	// DELETE endpoint to remove a recipe
	r.HandleFunc("/recipes", deleteRecipe).Methods("DELETE")

	// Start the local server
	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func getRecipes(w http.ResponseWriter, r *http.Request) {
	// A Query to retrieve all rows from the Recipes table
	rows, err := db.Query("SELECT id, title, description, ingredients, instructions, created_at, updated_at FROM Recipes")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through each row and Store the results in the array
	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &recipe.Ingredients, &recipe.Instructions, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define headers and Return recipes
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func deleteRecipe(w http.ResponseWriter, r *http.Request) {
	// Retrieve ID from query params
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Recipe ID", http.StatusBadRequest)
		return
	}

	// Delete a recipe from the database with the given ID
	result, err := db.Exec("DELETE FROM Recipes WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete Recipe", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
		return
	}
	
	if rowsAffected == 0 {
		http.Error(w, "No recipe found with the given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}