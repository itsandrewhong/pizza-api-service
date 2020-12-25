package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize a connection with the DB and initialize the router
func (a *App) Initialize() {
	connString := getConnString()

	var err error
	a.DB, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run the application
func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

// Initialize routes
func (a *App) initializeRoutes() {
	// Create a new customer
	a.Router.HandleFunc("/customer", a.createCustomer).Methods("POST")
	// Create a new order
	a.Router.HandleFunc("/order", a.createOrder).Methods("POST")
	// Get status of the order
	a.Router.HandleFunc("/order/{id:[0-9]+}", a.getStatus).Methods("GET")
}

// Helper: Get DB connection string from file
func getConnString() string {
	conString, err := ioutil.ReadFile("cstrings.config")

	if err != nil {
		log.Fatal(err)
	}

	return string(conString)
}

// Helper: Handle error message
func responseErrorHandler(w http.ResponseWriter, code int, message string) {
	responseWriter(w, code, map[string]string{"error": message})
}

// Helper: Write HTTP response in JSON format
func responseWriter(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Create a new customer give nthe HTTP request
func (a *App) createCustomer(w http.ResponseWriter, req *http.Request) {
	var p customer
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&p); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	defer req.Body.Close()

	if err := p.createCustomer(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWriter(w, http.StatusCreated, p)
}

// Create a new order given the HTTP request
func (a *App) createOrder(w http.ResponseWriter, req *http.Request) {

}

// Fetch the order status given the HTTP request
func (a *App) getStatus(w http.ResponseWriter, r *http.Request) {

}
