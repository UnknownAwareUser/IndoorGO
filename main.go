package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Jwt string `json:"token"`
}

var mySecretKey = []byte("veryverysecretkey")

func loginAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	ArrUser, _ := getQueryUser("SELECT login, pass FROM Administrator;")

	var t Token

	for _, user := range ArrUser {
		if u.Username == user.Username && u.Password == user.Password {
			t.Jwt, _ = GenTokenJWT()
			writeToken(u, t)
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func GenTokenJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix()
	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		log.Fatal()
	}

	return tokenString, nil
}

func checkAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySecretKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func writeToken(u User, t Token) {
	database, err := sql.Open("sqlite3", "api.db")
	defer database.Close()
	if err != nil {
		log.Fatal()
	}
	stmt, err := database.Prepare("UPDATE Administrator SET jwt = ? WHERE login = ?")
	defer stmt.Close()
	_, err = stmt.Exec(t.Jwt, u.Username)
	if err != nil {
		log.Fatal()
	}
	fmt.Println(t.Jwt)
}

func getQueryUser(call string) ([]User, []Token) {
	var Array = make([]User, 0)
	var ArrayToken = make([]Token, 0)
	database, err := sql.Open("sqlite3", "api.db")
	defer database.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := database.Query(call)
	if err != nil {
		return Array, ArrayToken
	}
	for rows.Next() {

		s := User{}
		sToken := Token{}
		err = rows.Scan(&s.Username, &s.Password)
		if err != nil {
			return Array, ArrayToken
		}
		Array = append(Array, s)
		ArrayToken = append(ArrayToken, sToken)
	}

	return Array, ArrayToken
}

//-------------------------------------------------------------------------------------------
//						Get data Shelf
//-------------------------------------------------------------------------------------------

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
	defer database.Close()
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
//						Get data Shelf
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
	defer database.Close()
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

//------------------------------------------------------------------------------------------------------------

// Data base create table and connect

func verDataBase() {
	database, err := sql.Open("sqlite3", "api.db")
	defer database.Close()
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

	statement, err1 = database.Prepare("CREATE TABLE IF NOT EXISTS Administrator (login TEXT NOT NULL, pass TEXT NOT NULL, jwt TEXT NOT NULL)")
	if err1 != nil {
		log.Fatal(err1)
	}

	statement.Exec()
	statement.Close()
}

//-------------------------------------------------------------------------------------------------------

// routing

func indoorAPI() *mux.Router {
	rout := mux.NewRouter()
	rout.HandleFunc("/api/v1/indoorforrest/list/all/shelfdata", getLastFive).Methods("GET")
	rout.HandleFunc("/api/v1/indoorforrest/list/all/shelfdata/{number:[1-5]}", getLastOne).Methods("GET")
	rout.HandleFunc("/api/v1/indoorforrest/list/all/workroom", getWorkRoom).Methods("GET")

	rout.Handle("/api/v1/admin/indoorforrest/list/all/shelfdata", checkAuthorized(getLastFive)).Methods("GET")
	rout.Handle("/api/v1/admin/indoorforrest/list/all/shelfdata/{number:[1-5]}", checkAuthorized(getLastOne)).Methods("GET")
	rout.Handle("/api/v1/admin/indoorforrest/list/all/workroom", checkAuthorized(getWorkRoom)).Methods("GET")

	// Authorization
	rout.HandleFunc("/api/v1/admin/login", loginAdmin).Methods("POST")
	return rout
}

func main() {
	verDataBase()
	fmt.Println("Start serv")
	log.Fatal(http.ListenAndServe("localhost:8080", indoorAPI()))

}
