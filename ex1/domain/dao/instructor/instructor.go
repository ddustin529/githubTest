package instructor

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

type instructor struct {
	Fn  string    `json:"Fn"`
	Ln  string    `json:"Ln"`
	IID uuid.UUID `json:"IID"`
}

func GetInstructors(w http.ResponseWriter, r *http.Request) {
	db = context.GetMySQLDB()
	defer db.Close()
	instructorList := []instructor{}
	rows, err := db.Query("select * from instructors")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var i instructor
		rows.Scan(&i.Fn, &i.Ln, &i.IID)
		instructorList = append(instructorList, i)
	}
	json.NewEncoder(w).Encode(instructorList)
}

func AddInstructors(w http.ResponseWriter, r *http.Request) {
	db := context.GetMySQLDB()
	defer db.Close()

	var i instructor
	i.IID = (uuid.New())
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO instructors(fn, ln, instructorID) VALUES (?, ?, ?)", i.Fn, i.Ln, i.IID)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(i)
}
