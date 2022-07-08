package cmd

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type house struct {
	address         string
	houseNumber     uint
	apartmentsCount uint
	floorsCount     uint
}

func dbConnection() (db *sql.DB) {
	dbDriver := "postgres"
	dbUser := "postgres"
	dbPass := "5gq5qe95"
	dbName := "apartmentdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	selDB, err := db.Query("SELECT * FROM houses ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	hse := house{}
	houses := []house{}
	for selDB.Next() {
		var houseNumber, apartmentsCount, floorsCount uint
		var address string
		err = selDB.Scan(&address, &houseNumber, &apartmentsCount, &floorsCount)
		if err != nil {
			panic(err.Error())
		}
		hse.address = address
		hse.houseNumber = houseNumber
		hse.apartmentsCount = apartmentsCount
		hse.floorsCount = floorsCount

		houses = append(houses, hse)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	nAddress := r.URL.Query().Get("address")
	selDB, err := db.Query("SELECT * FROM houses WHERE address=?", nAddress)
	if err != nil {
		panic(err.Error())
	}
	hse := house{}
	for selDB.Next() {
		var houseNumber, apartmentsCount, floorsCount uint
		var address string
		err = selDB.Scan(&address, &houseNumber, &apartmentsCount, &floorsCount)
		if err != nil {
			panic(err.Error())
		}
		hse.address = address
		hse.houseNumber = houseNumber
		hse.apartmentsCount = apartmentsCount
		hse.floorsCount = floorsCount
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	nAddress := r.URL.Query().Get("address")
	selDB, err := db.Query("SELECT * FROM houses WHERE address=?", nAddress)
	if err != nil {
		panic(err.Error())
	}
	hse := house{}
	for selDB.Next() {
		var houseNumber, apartmentsCount, floorsCount uint
		var address string
		err = selDB.Scan(&address, &houseNumber, &apartmentsCount, &floorsCount)
		if err != nil {
			panic(err.Error())
		}
		hse.address = address
		hse.houseNumber = houseNumber
		hse.apartmentsCount = apartmentsCount
		hse.floorsCount = floorsCount
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	if r.Method == "POST" {
		address := r.FormValue("address")
		houseNumber := r.FormValue("housenumber")
		apartmentsCount := r.FormValue("apartmentscount")
		floorsCount := r.FormValue("floorscount")
		insForm, err := db.Prepare("INSERT INTO houses(address, housenumber, apartmentscount, floorscount)" +
			" VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(address, houseNumber, apartmentsCount, floorsCount)
		log.Println("INSERT: Address: " + address + " | HouseNumber: " + houseNumber +
			"| ApartmentsCount: " + apartmentsCount + "| FloorsCount : " + floorsCount)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	if r.Method == "POST" {
		address := r.FormValue("address")
		houseNumber := r.FormValue("housenumber")
		apartmentsCount := r.FormValue("apartmentscount")
		floorsCount := r.FormValue("floorscount")
		insForm, err := db.Prepare("UPDATE houses SET housenumber=?, apartmentscount=?, floorscount =? WHERE address=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(address, houseNumber, apartmentsCount, floorsCount)
		log.Println("UPDATE: Address: " + address + " | HouseNumber: " + houseNumber +
			"| ApartmentsCount: " + apartmentsCount + "| FloorsCount : " + floorsCount)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	hse := r.URL.Query().Get("address")
	delForm, err := db.Prepare("DELETE FROM houses WHERE address=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(hse)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
