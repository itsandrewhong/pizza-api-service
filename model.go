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

// Takes firstName, lastName, and customerPhoneNumber, and
// creates a new row to 'CUSTOMERS' table with provided customer information, and returns the customerId
func (c *customer) createCustomer(db *sql.DB, hashedPassword string) error {
	// Calls the Stored Procedure and captures the customer id

	err := db.QueryRow("CALL PAS_SP_CREATE_CUSTOMER($1, $2, $3, $4, $5)", c.FirstName, c.LastName, c.CustomerPhoneNumber, c.Username, hashedPassword).Scan(&c.CustomerID)

	// err := db.QueryRow("INSERT INTO CUSTOMERS VALUES (DEFAULT, $1, $2, $3, $4, $5, FALSE) RETURNING customerId;", c.FirstName, c.LastName, c.CustomerPhoneNumber, c.Username, hashedPassword).Scan(&c.CustomerID)

	if err != nil {
		return err
	}

	return nil
}

// Takes pizzaId, and customerPhoneNumber, and
// creates a new row to 'ORDERS' table with provided order information, and returns the orderId
func (o *order) createOrder(db *sql.DB) error {
	// Calls the Stored Procedure and captures the order id
	err := db.QueryRow("CALL PAS_SP_CREATE_ORDER($1, $2)", o.PizzaID, o.CustomerPhoneNumber).Scan(&o.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// Takes in the orderId and returns the order status from 'ORDERS' table, and returns the orderStatus
func (s *status) getStatus(db *sql.DB, orderID int) error {
	// Calls the Stored Procedure 'PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER' and captures the order status
	return db.QueryRow("CALL PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER($1)", orderID).Scan(&s.StatusName)
}

// Cancels an order - Updates a statusId to '5' in 'ORDERS' table
func (s *status) cancelOrder(db *sql.DB, orderID int) error {
	// Calls the Stored Procedure 'PAS_SP_CANCEL_ORDER'
	return db.QueryRow("CALL PAS_SP_CANCEL_ORDER($1)", orderID).Scan(&s.StatusName)
}

// Retrieves list of orders by specific phone number
// Takes in the customerPhoneNumber and returns the list of orders by specific phone number
func (o *order) getOrders(db *sql.DB) ([]order, error) {
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

// Retrieves the list of available pizzas
func (p *pizza) getAvailablePizzas(db *sql.DB) ([]pizza, error) {
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

// Retrieves the list of status code (used by the store employees)
func (s *status) getStatusCode(db *sql.DB) ([]status, error) {
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

// Updates an OrderStatus for the specific order (used by the store employees) and returns the order status.
func (o *order) updateOrderStatus(db *sql.DB) error {
	// Calls the Stored Procedure 'PAS_SP_UPDATE_ORDER_STATUS'
	return db.QueryRow("CALL PAS_SP_UPDATE_ORDER_STATUS($1, $2)", o.OrderID, o.OrderStatus).Scan(&o.OrderStatus)
}

// Retrieves the hashed password given the username
func (c *customer) getCustomerPassword(db *sql.DB) error {

	// err := db.QueryRow("select password from CUSTOMERS where username=$1", c.Username).Scan(&c.Password)
	err := db.QueryRow("CALL PAS_SP_GET_CUSTOMER_PASSWORD($1)", c.Username).Scan(&c.Password)

	if err != nil {
		return err
	}

	return nil
}
