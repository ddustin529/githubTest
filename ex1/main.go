package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"

	"main/domain/dao/course"
	"main/domain/dao/instructor"
	"main/domain/dao/student"
)

var db *sql.DB

func main() {
	db, err := sql.Open("mysql", "dustin:dad1022@tcp(127.0.0.1:3306)/student")
	if err != nil {
		fmt.Println("open panic")
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("successfully connected")

	router := mux.NewRouter()
	router.HandleFunc("/students", student.GetStudents).Methods("GET")
	router.HandleFunc("/students/{id:[a-f0-9\\-]+}", student.GetStudentByID).Methods("GET")
	router.HandleFunc("/instructors", instructor.GetInstructors).Methods("GET")
	router.HandleFunc("/courses", course.GetCourses).Methods("GET")
	router.HandleFunc("/students", student.AddStudents).Methods("POST")
	router.HandleFunc("/instructors", instructor.AddInstructors).Methods("POST")
	router.HandleFunc("/courses", course.AddCourses).Methods("POST")
	router.HandleFunc("/students/{id:[a-f0-9\\-]+}", student.UpdateStudents).Methods("PUT")
	router.HandleFunc("/students/{id:[a-f0-9\\-]+}", student.DeleteStudents).Methods("DELETE")

	http.ListenAndServe(":8080", router)

}
