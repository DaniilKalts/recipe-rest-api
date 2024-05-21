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
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Ingredient struct {
	Ingredient string `json:"ingredient"`
	Quantity   string `json:"quantity"`
}

type Recipe struct {
	Id           int          `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
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

	// POST endpoint to create a new recipe
	r.HandleFunc("/recipes", createRecipe).Methods("POST")

	// UPDATE endpoint to update an existing recipe
	r.HandleFunc("/recipes", updateRecipe).Methods("PUT")

	// Start the local server
	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func getRecipes(w http.ResponseWriter, r *http.Request) {
	// A Query to retrieve all recipes from the Recipes table
	rows, err := db.Query("SELECT id, title, description, ingredients, instructions, created_at, updated_at FROM Recipes")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the recipes
	var recipes []Recipe

	// Iterate through the rows
	for rows.Next() {
		var recipe Recipe

		// Scan the values from the row into variables
		var ingredientsJSON []byte // Create a byte slice to hold the JSON data
		var instructions []string  // Create a slice of strings to hold the instructions array data
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description, &ingredientsJSON, pq.Array(&instructions), &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the ingredients JSON into a slice of Ingredient structs
		if err := json.Unmarshal(ingredientsJSON, &recipe.Ingredients); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Assign the instructions directly to the recipe
		recipe.Instructions = instructions

		// Append the recipe to the slice
		recipes = append(recipes, recipe)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode recipes slice as JSON and send it in the response
	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully fetched all recipes:", recipes)
}

func deleteRecipe(w http.ResponseWriter, r *http.Request) {
	// Retrieve ID from query params
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Recipe ID", http.StatusBadRequest)
		log.Println("Missing Recipe ID")
		return
	}

	// Delete a recipe from the database with the given ID
	result, err := db.Exec("DELETE FROM Recipes WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete Recipe", http.StatusInternalServerError)
		log.Println("Failed to delete Recipe", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
		log.Println("Failed to get affected rows", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No recipe found with the given ID", http.StatusNotFound)
		log.Println("No recipe found with the given ID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Successfully deleted recipe with ID:", id)
}

func createRecipe(w http.ResponseWriter, r *http.Request) {
	// Decode the request body to get a new recipe data
	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Failed to decode request payload:", err)
		return
	}
	defer r.Body.Close()

	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	// Convert Ingredients to JSON before inserting into the database
	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		http.Error(w, "Failed to marshal ingredients", http.StatusInternalServerError)
		log.Println("Failed to marshal ingredients:", err)
		return
	}

	// Convert Instructions to an array format before inserting into the database
	instructionsArray := pq.Array(recipe.Instructions)

	// Create a new recipe query
	query := `INSERT INTO Recipes (title, description, ingredients, instructions, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = db.QueryRow(query, recipe.Title, recipe.Description, ingredientsJSON, instructionsArray, recipe.CreatedAt, recipe.UpdatedAt).Scan(&recipe.Id)
	if err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		log.Println("Failed to insert a new Recipe", err)
		return
	}

	// Set response headers and encode a new recipe as JSON
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
	log.Println("Successfully created a new recipe with ID:", recipe.Id)
}

func updateRecipe(w http.ResponseWriter, r *http.Request) {
	// Retrieve ID from query params
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Recipe ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated recipe data
	var updatedRecipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&updatedRecipe); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Failed to decode request payload:", err)
		return
	}
	defer r.Body.Close()

	// Convert Ingredients to JSON before updating in the database
	ingredientsJSON, err := json.Marshal(updatedRecipe.Ingredients)
	if err != nil {
		http.Error(w, "Failed to marshal ingredients", http.StatusInternalServerError)
		log.Println("Failed to marshal ingredients:", err)
		return
	}

	// Convert Instructions to an array format before updating in the database
	instructionsArray := pq.Array(updatedRecipe.Instructions)

	// Execute the update query
	query := `UPDATE Recipes SET title = $1, description = $2, ingredients = $3, instructions = $4, updated_at = $5 WHERE id = $6`
	result, err := db.Exec(query, updatedRecipe.Title, updatedRecipe.Description, ingredientsJSON, instructionsArray, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		log.Println("Failed to update recipe:", err)
		return
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		log.Println("Failed to retrieve affected rows:", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No recipe found with the given ID", http.StatusNotFound)
		log.Println("No recipe found with the given ID:", id)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully updated recipe with ID:", id)
}
