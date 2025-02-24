package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Kết nối Database
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "Main:123456@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println(" Ket noi MySQL thanh cong!")
	return db, nil
}
func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf(" Loi ket noi database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			addUser(w, r, db)
		case "PUT":
			updateUser(w, r, db)
		default:
			http.Error(w, "Phuong thuc khong ho tro", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			http.Error(w, "Thieu ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case "GET":
			getUserByID(w, id, db)
		case "DELETE":
			deleteUser(w, id, db)
		default:
			http.Error(w, "Phuong thuc khong ho tro", http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// API: Thêm User
func addUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user struct{ Name, Password string }
	if json.NewDecoder(r.Body).Decode(&user) != nil {
		http.Error(w, "Du lieu khong hop le", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		http.Error(w, "Loi khi them user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "Them user thanh cong!")
}

// API: Cập nhật User
func updateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user struct{ Name, Password string }
	if json.NewDecoder(r.Body).Decode(&user) != nil {
		http.Error(w, "Du lieu khong hop le", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE users SET password = ? WHERE name = ?", user.Password, user.Name)
	if err != nil {
		http.Error(w, "Loi khi cap nhat user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "Cap nhat user thanh cong!")
}

// API: Xóa User
func deleteUser(w http.ResponseWriter, id string, db *sql.DB) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID khong hop le", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM users WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, "Loi khi xoa user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "Xoa user thanh cong!")
}

// API: Lấy User theo ID
func getUserByID(w http.ResponseWriter, id string, db *sql.DB) {
	var user struct{ Name, Password string }
	err := db.QueryRow("SELECT name, password FROM users WHERE id = ?", id).Scan(&user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User khong ton tai", http.StatusNotFound)
		} else {
			http.Error(w, "Loi khi lay user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Trả về JSON response
func jsonResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
