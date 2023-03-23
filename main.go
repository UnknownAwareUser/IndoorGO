package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type shelf struct {
	Number_shelf   int `json:"number_Shelf"`
	Humindt_flower int `json:"humindt_Flower"`
	Lighting_shelf int `json:"lighting_Shelf"`
	Updated        int `json:"updated"`
	// (id INTEGER PRIMARY KEY, number INTEGER, luminosity INTEGER, humidity INTEGER, updated INTEGER)
}

type workroom struct {
	Id              int `json:"id"`
	Air_temperature int `json:"air_temperature"`
	Air_humidity    int `json:"air_humidity"`
	Air_pressure    int `json:"air_pressure"`
	Updated         int `json:"updated"`
	// (id INTEGER PRIMARY KEY, temperature INTEGER, humidity INTEGER, pressure INTEGER, updated INTEGER)
}

func getLastFive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	ArrShelf := getQueryShelf("SELECT MAX(updated), number, luminosity, humidity FROM shelf GROUP BY number")
	json.NewEncoder(w).Encode(ArrShelf)
}

func getLastOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	s := mux.Vars(r)["number"]

	ArrShelf := getQueryShelf("SELECT MAX(updated), number, luminosity, humidity FROM shelf WHERE number ==" + s)
	json.NewEncoder(w).Encode(ArrShelf)
	// INSERT INTO stand (number, luminosity, humidity, updated)
	// VALUES (1, 2, 7, 1);
}

func getQueryShelf(call string) []shelf {
	var Array = make([]shelf, 0)
	database, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := database.Query(call)
	if err != nil {
		return Array
	}
	for rows.Next() {
		s := shelf{}
		err = rows.Scan(&s.Updated, &s.Number_shelf, &s.Lighting_shelf, &s.Humindt_flower)
		if err != nil {
			return Array
		}
		Array = append(Array, s)
	}

	return Array
}

//-------------------------------------------------------------------------------------------

func getWorkRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	Arr := getQueryWorkRoom("SELECT MAX(updated), temperature, humidity, pressure, id FROM stand")
	json.NewEncoder(w).Encode(Arr)
	// INSERT INTO stand (temperature, humidity, pressure, updated)
	// VALUES (3, 2, 7, 1);
}

func getQueryWorkRoom(call string) []workroom {
	var Array = make([]workroom, 0)
	database, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := database.Query(call)
	if err != nil {
		return Array
	}
	for rows.Next() {
		s := workroom{}
		err = rows.Scan(&s.Updated, &s.Air_temperature, &s.Air_humidity, &s.Air_pressure, &s.Id)
		if err != nil {
			return Array
		}
		Array = append(Array, s)
	}

	return Array
}

func indoorAPI() *mux.Router {
	rout := mux.NewRouter()
	rout.HandleFunc("/api/v1/indoorforrest/list/all/shelfdata/", getLastFive)
	rout.HandleFunc("/api/v1/indoorforrest/list/all/shelfdata/{number:[1-5]}/", getLastOne)
	rout.HandleFunc("/api/v1/indoorforrest/list/all/workroom/", getWorkRoom)
	return rout
}

func verDataBase() {
	database, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err1 := database.Prepare("CREATE TABLE IF NOT EXISTS stand (id INTEGER PRIMARY KEY, temperature INTEGER, humidity INTEGER, pressure INTEGER, updated INTEGER)")
	if err1 != nil {
		log.Fatal(err1)
	}
	statement.Exec()

	statement, err1 = database.Prepare("CREATE TABLE IF NOT EXISTS shelf (number INTEGER, luminosity INTEGER, humidity INTEGER, updated INTEGER, PRIMARY KEY(number, updated))")
	if err1 != nil {
		log.Fatal(err1)
	}
	statement.Exec()
}

func main() {
	verDataBase()
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", indoorAPI()))
}
