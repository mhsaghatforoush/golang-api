package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "987654321"
	dbname   = "api"
)

var db *sql.DB

func init() {
	var err error
	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create the users table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		phone_number VARCHAR(20) NOT NULL
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate phone number (add your validation logic here)

	// Insert the new user into the database
	var userID int
	err = db.QueryRow("INSERT INTO users (phone_number) VALUES ($1) RETURNING user_id", user.PhoneNumber).Scan(&userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"message": "User registered successfully", "user_id": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", registerHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", r))
}
