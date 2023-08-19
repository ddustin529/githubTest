package student

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"main/service/context"
)

var db *sql.DB

type student struct {
	Fn  string    `json:"fn"`
	Ln  string    `json:"ln"`
	SID uuid.UUID `json:"sid"`
}

// GETTING
func GetStudents(w http.ResponseWriter, r *http.Request) {
	db = context.GetMySQLDB()
	defer db.Close()
	studentList := []student{}
	rows, err := db.Query("select * from students")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var s student
		rows.Scan(&s.Fn, &s.Ln, &s.SID)
		studentList = append(studentList, s)
	}
	json.NewEncoder(w).Encode(studentList)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()

	vars := mux.Vars(r)
	studentID := vars["id"]

	var s student
	err := db.QueryRow("SELECT fn, ln, studentID FROM students WHERE studentID = ?", studentID).
		Scan(&s.Fn, &s.Ln, &s.SID)
	if err != nil {
		http.Error(w, "Error fetching student data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// POSTING
func AddStudents(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()

	var s student
	s.SID = (uuid.New())
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO students(fn, ln, studentID) VALUES (?, ?, ?)", s.Fn, s.Ln, s.SID)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(s)
}

// PUTTING
func UpdateStudents(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()
	var s student
	json.NewDecoder(r.Body).Decode(&s)
	result, err := db.Exec("update students set fn=?, ln=? where studentID=?", s.Fn, s.Ln, s.SID)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{error:Record is not updated}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}
}

// DELETING
func DeleteStudents(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()
	vars := mux.Vars(r)
	studentID := vars["id"]
	result, err := db.Exec("delete from students where studentID=?", studentID)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{error:Record is not deleted}")
		} else {
			json.NewEncoder(w).Encode("{Record is deleted}")
		}
	}
}
