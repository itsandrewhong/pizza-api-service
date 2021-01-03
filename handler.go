package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

	// Validate customer phone number
	if err := validateCustomerPhoneNumber(c); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write customer data to DB
	if err := c.createCustomer(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	payload := map[string]string{
		"customerPhoneNumber": c.CustomerPhoneNumber,
	}

	// Write HTTP response
	responseWriter(w, http.StatusCreated, payload)
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

	// Validate customer phone number
	if err := validateCustomerPhoneNumber(o); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write order data to DB
	if err := o.createOrder(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	payload := map[string]int{
		"orderId": o.OrderID,
	}

	// Write HTTP response
	responseWriter(w, http.StatusCreated, payload)
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

	// Get the current order status from DB
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

	// Create a HTTP Response payload
	payload := map[string]string{
		"orderStatus": s.StatusName,
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, payload)
}

// Handler to cancel an order.
func (a *App) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Create route variable and retrieve 'orderId' from a Request URL
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		responseErrorHandler(w, http.StatusBadRequest, "Invalid order ID")
		return
	}
	var s status

	// Update a row in DB
	if err := s.cancelOrder(a.DB, orderID); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	payload := map[string]string{
		"orderStatus": s.StatusName,
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, payload)
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

	// Validate customer phone number
	if err := validateCustomerPhoneNumber(o); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get order data from DB
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

	// Get the list of available pizzas from DB
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

	// Get order status from DB
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

	// Update a row in DB
	if err := o.updateOrderStatus(a.DB); err != nil {
		responseErrorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a HTTP Response payload
	payload := map[string]interface{}{
		"orderId":     strconv.Itoa(o.OrderID),
		"orderStatus": fmt.Sprintf("%v", o.OrderStatus),
	}

	// Write HTTP response
	responseWriter(w, http.StatusOK, payload)
}
