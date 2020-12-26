package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
func (a *App) Run(addr string) {
	log.Println("listen on", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Initialize routes
func (a *App) initializeRoutes() {
	// Create a new customer
	a.Router.HandleFunc("/customer", a.createCustomerHandler).Methods("POST")
	// Create a new order
	a.Router.HandleFunc("/order", a.createOrderHandler).Methods("POST")
	// Get status of the order
	a.Router.HandleFunc("/order/{id:[0-9]+}", a.getStatusHandler).Methods("GET")
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

// Handler to create a new customer.
// Takes a request body in JSON format and uses 'createCustomer' to create a customer.
func (a *App) createCustomerHandler(w http.ResponseWriter, req *http.Request) {
	var c customer
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&c); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer req.Body.Close()

	if err := c.createCustomer(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWriter(w, http.StatusCreated, c)
}

// Handler to create a new order.
func (a *App) createOrderHandler(w http.ResponseWriter, req *http.Request) {
	var o order
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&o); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer req.Body.Close()

	if err := o.createOrder(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWriter(w, http.StatusCreated, o)
}

// Handler to fetch the order status.
// Retrieves the order id and returns the order status.
// If order is not found, respond with status code 404,
// If found, return the order status.
func (a *App) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	o := order{OrderID: orderID}
	if err := o.getStatus(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responseErrorHandler(w, http.StatusNotFound, "Order not found")
		default:
			responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responseWriter(w, http.StatusOK, o)
}
