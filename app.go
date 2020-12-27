package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize a connection with the DB and initialize the router
func (a *App) Initialize() {
	// Connect to the DB (Heroku Postgres)
	var err error
	a.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new router and initialize routes
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
	a.Router.HandleFunc("/customer/add", a.createCustomerHandler).Methods("POST")
	// Create a new order
	a.Router.HandleFunc("/order/add", a.createOrderHandler).Methods("POST")
	// Get status of the order
	a.Router.HandleFunc("/order/show/{orderId:[0-9]+}", a.getStatusHandler).Methods("GET")
	// Cancel an order
	a.Router.HandleFunc("/order/update/{orderId:[0-9]+}", a.cancelOrderHandler).Methods("PUT")
	// Get list of orders by specific phone number
	a.Router.HandleFunc("/order/show", a.getOrdersHandler).Methods("GET")
	// Get the list of available pizzas
	a.Router.HandleFunc("/pizza/show", a.getAvailablePizzasHandler).Methods("GET")

	// Get the list of order status
	a.Router.HandleFunc("/order/codes", a.getStatusCodeHandler).Methods("GET")

	// Update the order status
	a.Router.HandleFunc("/order/update", a.updateOrderStatusHandler).Methods("PUT")
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

// CustomerPhoneNumber cannot be empty, and must be a string consisting of ten digits
func validateCustomerPhoneNumber(i interface{}) error {
	switch v := i.(type) {
	case order:
		return validation.ValidateStruct(&v, validation.Field(&v.CustomerPhoneNumber, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{10}$"))))
	case customer:
		return validation.ValidateStruct(&v, validation.Field(&v.CustomerPhoneNumber, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{10}$"))))
	}

	return nil
}

// Handler to create a new customer.
// Takes a request body in JSON format and uses 'createCustomer' to create a customer.
func (a *App) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var c customer
	decoder := json.NewDecoder(r.Body)

	// Decode the HTTP Body Data
	if err := decoder.Decode(&c); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := validateCustomerPhoneNumber(c); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write to DB
	if err := c.createCustomer(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusCreated, c)
}

// Handler to create a new order.
func (a *App) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var o order
	decoder := json.NewDecoder(r.Body)

	// Decode the HTTP Body Data
	if err := decoder.Decode(&o); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := validateCustomerPhoneNumber(o); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write to DB
	if err := o.createOrder(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusCreated, o)
}

// Handler to fetch the order status.
// Retrieves the order id and returns the order status.
// If order is not found, respond with status code 404, if found, return the order status.
func (a *App) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Create route variable and retrieve 'orderId' from a Request URL
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Get order status from DB
	var s status
	if err := s.getStatus(a.DB, orderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			responseErrorHandler(w, http.StatusNotFound, "Order not found")
		default:
			responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, s.StatusName)
}

// Handler to cancel an order
func (a *App) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Create route variable and retrieve 'orderId' from a Request URL
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// var o order
	// o.OrderID = orderID

	var s status

	// Write to DB (Update a row)
	if err := s.cancelOrder(a.DB, orderID); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	// payload := map[string]string{
	// 	"orderId":    strconv.Itoa(orderID),
	// 	"statusName": s.StatusName,
	// }

	// Write HTTP response
	responseWriter(w, http.StatusOK, s.StatusName)
}

// Handler to fetch orders given the customer phone number
func (a *App) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var o order
	decoder := json.NewDecoder(r.Body)

	// Decode the HTTP Body Data
	if err := decoder.Decode(&o); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := validateCustomerPhoneNumber(o); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get data from DB
	orders, err := o.getOrders(a.DB)
	if err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, orders)
}

// Handler to fetch the list of available pizzas
func (a *App) getAvailablePizzasHandler(w http.ResponseWriter, r *http.Request) {
	var p pizza

	// Get data from DB
	orders, err := p.getAvailablePizzas(a.DB)
	if err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, orders)
}

// Handler to fetch the list of order status
// Used by store employees
func (a *App) getStatusCodeHandler(w http.ResponseWriter, r *http.Request) {
	var s status

	// Get data from DB
	orders, err := s.getStatusCode(a.DB)
	if err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, orders)
}

// Handler to update an order status
// Receives HTTP Body data
// Used by store employees
func (a *App) updateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	var o order
	decoder := json.NewDecoder(r.Body)

	// Decode the HTTP Body Data
	if err := decoder.Decode(&o); err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	defer r.Body.Close()

	// Write to DB (Update a row)
	if err := o.updateOrderStatus(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	// Payload contains the orderId and orderStatus
	payload := map[string]interface{}{
		"orderId":     strconv.Itoa(o.OrderID),
		"orderStatus": fmt.Sprintf("%v", o.OrderStatus),
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, payload)
}
