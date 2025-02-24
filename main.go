package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Connect to Database
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "Main:123456@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println(" Successfully connected to MySQL!")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf(" Database connection error: %v", err)
	}
	defer db.Close()
	http.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		AddUser(w, r, db)
	})
	http.HandleFunc("PUT /users", func(w http.ResponseWriter, r *http.Request) {
		UpdateUser(w, r, db)
	})
	http.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		DeleteUser(w, id, db)
	})
	http.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		GetAllUsers(w, db)
	})
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		GetUserByID(w, id, db)
	})
	http.ListenAndServe(":8080", nil)
}

// API: Create User
func AddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user struct{ Name, Password string }
	if json.NewDecoder(r.Body).Decode(&user) != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "User added successfully!")
}

// API: Update User
func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user struct{ ID, Name, Password string }
	if json.NewDecoder(r.Body).Decode(&user) != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(user.ID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE users SET name = ?, password = ? WHERE id = ?", user.Name, user.Password, idInt)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "User updated successfully!")
}

// API: Get All Users
func GetAllUsers(w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, password string
		if err := rows.Scan(&id, &name, &password); err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]interface{}{"id": id, "name": name, "password": password})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// API: Get User by ID
func GetUserByID(w http.ResponseWriter, id string, db *sql.DB) {
	var user struct{ Name, Password string }
	err := db.QueryRow("SELECT name, password FROM users WHERE id = ?", id).Scan(&user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// API: Delete User by ID
func DeleteUser(w http.ResponseWriter, id string, db *sql.DB) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, "User deleted successfully")
}

func jsonResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
