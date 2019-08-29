package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var dbName = readEnvFile("dbName")
var dbPort = readEnvFile("dbPort")
var dbUser = readEnvFile("dbUser")
var dbPasswd = readEnvFile("dbPasswd")
var dbHost = readEnvFile("dbHost")

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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbHost, dbPort, dbName))

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

func readEnvFile(varName string) string {
	file, err := os.Open(".env")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), varName) {
			return strings.Split(scanner.Text(), "=")[1]
		}
	}
	return ""
}

func main() {
	serverPort := 8000
	r := mux.NewRouter()

	r.HandleFunc("/health-check", getHealthStatus).Methods("GET")
	r.HandleFunc("/api/add-user", dbInsert).Methods("POST")

	log.Println("Starting web server on port " + strconv.Itoa(serverPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverPort), r))
}
