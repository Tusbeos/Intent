package main

import (
	models "Intent/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Ket noi Database
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "Main:123456@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MySQL!")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
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
	http.HandleFunc("GET /users/", func(w http.ResponseWriter, r *http.Request) {
		GetListUsers(w, r, db)
	})
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		GetUserByID(w, id, db)
	})
	http.ListenAndServe(":8080", nil)
}

// API: Them nguoi dung
func AddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, 400, "Invalid data", nil)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		jsonResponse(w, 500, "Error adding user", nil)
		return
	}
	jsonResponse(w, 0, "User added successfully!", nil)
}

// API: Cap nhat nguoi dung
func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, 400, "Invalid data", nil)
		return
	}

	result, err := db.Exec("UPDATE users SET name = ?, password = ? WHERE id = ?", user.Name, user.Password, user.ID)
	if err != nil {
		jsonResponse(w, 500, "Error updating user", nil)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonResponse(w, 404, "User not found", nil)
		return
	}

	jsonResponse(w, 0, "User updated successfully!", nil)
}

// API: Lay danh sach nguoi dung
func GetListUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Lay gia tri page và limit
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	//Tinh total
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		jsonResponseMeta(w, 500, "Error counting users", nil, models.Meta{})
		return
	}

	// Lay danh sach phan trang
	offset := (page - 1) * limit
	rows, err := db.Query("SELECT id, name, password FROM users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		jsonResponseMeta(w, 500, "Error retrieving user list", nil, models.Meta{})
		return
	}
	defer rows.Close()

	// Duyet qua ket qua và luu vao danh sach user
	var users []models.Users
	for rows.Next() {
		var user models.Users
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			jsonResponseMeta(w, 500, "Error reading data", nil, models.Meta{})
			return
		}
		users = append(users, user)
	}
	// Tao object meta
	meta := models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	// Tra ve JSON response
	jsonResponseMeta(w, 0, "User list retrieved successfully", users, meta)
}

// Hàm tra ve JSON response
func jsonResponseMeta(w http.ResponseWriter, errorCode int, message string, data interface{}, meta models.Meta) {
	w.Header().Set("Content-Type", "application/json")
	response := models.ResponseMeta{
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
		Meta:      meta,
	}
	json.NewEncoder(w).Encode(response)
}

// API: Lay thong tin nguoi dung theo ID
func GetUserByID(w http.ResponseWriter, id string, db *sql.DB) {
	var user models.Users
	err := db.QueryRow("SELECT id, name, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, 404, "User not found", nil)
		} else {
			jsonResponse(w, 500, "Error retrieving user", nil)
		}
		return
	}
	jsonResponse(w, 0, "User retrieved successfully", user)
}

// API: Xoa nguoi dung theo ID
func DeleteUser(w http.ResponseWriter, id string, db *sql.DB) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		jsonResponse(w, 400, "Invalid ID", nil)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = ?", idInt)
	if err != nil {
		jsonResponse(w, 500, "Error deleting user", nil)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonResponse(w, 404, "User not found", nil)
		return
	}

	jsonResponse(w, 0, "User deleted successfully", nil)
}

// Ham JSON response
func jsonResponse(w http.ResponseWriter, errorCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
	}
	json.NewEncoder(w).Encode(response)
}
