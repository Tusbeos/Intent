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
	http.HandleFunc("/POST/users", func(w http.ResponseWriter, r *http.Request) {
		AddUser(w, r, db)
	})
	http.HandleFunc("/PUT/users", func(w http.ResponseWriter, r *http.Request) {
		UpdateUser(w, r, db)
	})
	http.HandleFunc("/DELETE/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/DELETE/users/")
		DeleteUser(w, id, db)
	})
	http.HandleFunc("/GET/users", func(w http.ResponseWriter, r *http.Request) {
		GetAllUsers(w, db)
	})
	http.HandleFunc("/GET/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/GET/users/")
		GetUserByID(w, id, db)
	})
	http.ListenAndServe(":8080", nil)
}

// API: Tạo User
func AddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user struct{ ID, Name, Password string }
	if json.NewDecoder(r.Body).Decode(&user) != nil {
		http.Error(w, "Du lieu khong hop le", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(user.ID)
	if err != nil {
		http.Error(w, "ID khong hop le", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE users SET name = ?, password = ? WHERE id = ?", user.Name, user.Password, idInt)
	if err != nil {
		http.Error(w, "Loi khi cap nhat user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, "Cap nhat user thanh cong!")
}

// API: Lay danh sach User
func GetAllUsers(w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		http.Error(w, "Loi khi lay danh sach users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, password string
		if err := rows.Scan(&id, &name, &password); err != nil {
			http.Error(w, "Loi khi doc du lieu", http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]interface{}{"id": id, "name": name, "password": password})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// API: Lấy User theo ID
func GetUserByID(w http.ResponseWriter, id string, db *sql.DB) {
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

// API: Xóa User theo ID
func DeleteUser(w http.ResponseWriter, id string, db *sql.DB) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID khong hop le", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, "Lỗi khi xóa user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User khong ton tai", http.StatusNotFound)
		return
	}

	jsonResponse(w, "Xoa user thanh cong")
}

func jsonResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
