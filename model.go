package main

import (
	"database/sql"
	"time"
)

// Create a struct that holds the 'customer' information
type customer struct {
	CustomerID          int    `json:"customerId"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	Username            string `json:"username"`
	Password            string `json:"password"`
}

// Create a struct that holds the 'order' information
type order struct {
	OrderID             int         `json:"orderId"`
	PizzaID             int         `json:"pizzaId"`
	OrderTime           time.Time   `json:"orderTime"`
	CustomerPhoneNumber string      `json:"customerPhoneNumber"`
	OrderStatus         interface{} `json:"orderStatus"`
	TotalPrice          float64     `json:"totalPrice"`
}

// Create a struct that holds the 'pizza' information
type pizza struct {
	PizzaID    int     `json:"pizzaId"`
	PizzaName  string  `json:"pizzaName"`
	PizzaPrice float64 `json:"pizzaPrice"`
}

// Create a struct that holds the 'status' information
type status struct {
	StatusID   int    `json:"statusId"`
	StatusName string `json:"statusName"`
}

// Create a customer
// Takes firstName, lastName, and customerPhoneNumber
// Create a new row to 'CUSTOMERS' table with provided customer information, and returns the customerId
func (c *customer) createCustomer(db *sql.DB) error {
	// Calls the Stored Procedure and captures the customer id
	err := db.QueryRow("CALL PAS_SP_CREATE_CUSTOMER($1, $2, $3)", c.FirstName, c.LastName, c.CustomerPhoneNumber).Scan(&c.CustomerID)
	if err != nil {
		return err
	}

	return nil
}

// Create an order
// Takes pizzaId, and customerPhoneNumber
// Create a new row to 'ORDERS' table with provided order information, and returns the orderId
func (o *order) createOrder(db *sql.DB) error {
	// Calls the Stored Procedure and captures the order id
	err := db.QueryRow("CALL PAS_SP_CREATE_ORDER($1, $2)", o.PizzaID, o.CustomerPhoneNumber).Scan(&o.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// Get order status
// Takes in the orderId
// Retrieves the order status from 'ORDERS' table, and returns the orderStatus
func (s *status) getStatus(db *sql.DB, orderID int) error {
	// Calls the Stored Procedure 'PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER' and captures the order status
	return db.QueryRow("CALL PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER($1)", orderID).Scan(&s.StatusName)
}

// Cancel an order
// Updates a statusId in 'ORDERS' table
func (s *status) cancelOrder(db *sql.DB, orderID int) error {
	// Calls the Stored Procedure 'PAS_SP_CANCEL_ORDER'
	return db.QueryRow("CALL PAS_SP_CANCEL_ORDER($1)", orderID).Scan(&s.StatusName)
}

// Get list of orders by specific phone number
// Takes in the customer phone number and returns the list of orders by specific phone number
func (o *order) getOrders(db *sql.DB) ([]order, error) {
	// Run the query
	rows, err := db.Query(
		"SELECT o.customerPhoneNumber, o.orderId, o.orderTime, o.pizzaId, o.totalPrice, sc.statusName FROM ORDERS AS o INNER JOIN ORDER_STATUS_CODES AS sc ON o.statusId = sc.statusId WHERE customerPhoneNumber = $1 AND o.isDeleted = FALSE", o.CustomerPhoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a 'orders' list and append each resulting row to the 'orders' list
	orders := []order{}
	for rows.Next() {
		if err := rows.Scan(&o.CustomerPhoneNumber, &o.OrderID, &o.OrderTime, &o.PizzaID, &o.TotalPrice, &o.OrderStatus); err != nil {
			return nil, err
		}
		orders = append(orders, *o)
	}

	return orders, nil
}

// Get the list of available pizzas
// Returns the list of available pizzas
func (p *pizza) getAvailablePizzas(db *sql.DB) ([]pizza, error) {
	// Run the query
	rows, err := db.Query(
		"SELECT pizzaId, pizzaName, pizzaPrice FROM PIZZAS WHERE isDeleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a 'pizzas' list and append each resulting row to the 'pizzas' list
	pizzas := []pizza{}
	for rows.Next() {
		if err := rows.Scan(&p.PizzaID, &p.PizzaName, &p.PizzaPrice); err != nil {
			return nil, err
		}
		pizzas = append(pizzas, *p)
	}

	return pizzas, nil
}

// Get the list of status code (used by the store employees)
// Returns the list of status codes
func (s *status) getStatusCode(db *sql.DB) ([]status, error) {
	// Run the query
	rows, err := db.Query("SELECT statusId, statusName FROM ORDER_STATUS_CODES WHERE isDeleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a 'statuses' list and append each resulting row to the 'statuses' list
	statuses := []status{}
	for rows.Next() {
		if err := rows.Scan(&s.StatusID, &s.StatusName); err != nil {
			return nil, err
		}
		statuses = append(statuses, *s)
	}

	return statuses, nil
}

// Updates the order status (used by the store employees)
// Updates a row on 'ORDERS' table and returns the order status.
func (o *order) updateOrderStatus(db *sql.DB) error {
	// Calls the Stored Procedure 'PAS_SP_UPDATE_ORDER_STATUS'
	return db.QueryRow("CALL PAS_SP_UPDATE_ORDER_STATUS($1, $2)", o.OrderID, o.OrderStatus).Scan(&o.OrderStatus)
}
