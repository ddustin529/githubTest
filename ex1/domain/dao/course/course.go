package course

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"main/service/context"
)

var db *sql.DB

type course struct {
	CID uuid.UUID `json:"cid"`
	CN  string    `json:"cn"`
	IID uuid.UUID `json:"iid"`
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	db = context.GetMySQLDB()
	defer db.Close()
	courseList := []course{}
	rows, err := db.Query("select * from courses")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var c course
		rows.Scan(&c.CID, &c.CN, &c.IID)
		courseList = append(courseList, c)
	}
	json.NewEncoder(w).Encode(courseList)
}

func AddCourses(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()
	var c course
	c.CID = (uuid.New())
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO courses(courseID, title, instructorID) VALUES (?, ?, ?)", c.CID, c.CN, c.IID)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
}
