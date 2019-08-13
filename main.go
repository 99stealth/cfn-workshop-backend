package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func getHealthStatus(w http.ResponseWriter, r *http.Request) {
	var version = "0.0.1"
	var status = "healthy"
	var serviceName = "cfn-workshop"
	healthCheck := map[string]string{
		"Application version": version,
		"Application status":  status,
		"Service Name":        serviceName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthCheck)
}

func dbClient() *sql.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPasswd := "5ecurePa$$word"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "cfn_workshop"
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbHost, dbPort, dbName))

	if err != nil {
		panic(err.Error())
	}
	return db
}

func dbInsert(w http.ResponseWriter, r *http.Request) {
	db := dbClient()
	if r.Method == "POST" {
		timestamp := time.Now()
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		insForm, err := db.Prepare("INSERT INTO user(firstname, lastname, email, reg_date) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(firstname, lastname, email, timestamp.Format("2006-01-02 15:04:05"))
		log.Printf("INSERT -> First name: %s, Last name: %s, Email: %s", firstname, lastname, email)
	}
	defer db.Close()
}

func main() {
	serverPort := 8000
	r := mux.NewRouter()

	r.HandleFunc("/health-check", getHealthStatus).Methods("GET")
	r.HandleFunc("/api/add-user", dbInsert).Methods("POST")

	log.Println("Starting web server on port " + strconv.Itoa(serverPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverPort), r))
}
