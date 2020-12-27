package main

import (
	"database/sql"
	"time"
)

type customer struct {
	CustomerID          int    `json:"customerId"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
}

type order struct {
	OrderID             int       `json:"orderId"`
	PizzaID             int       `json:"pizzaId"`
	OrderTime           time.Time `json:"orderTime"`
	CustomerPhoneNumber string    `json:"customerPhoneNumber"`
	OrderStatus         string    `json:"orderStatus"`
	TotalPrice          float64   `json:"totalPrice"`
}

type pizza struct {
	PizzaID    int     `json:"pizzaId"`
	PizzaName  string  `json:"pizzaName"`
	PizzaPrice float64 `json:"pizzaPrice"`
}

// Create a customer
func (c *customer) createCustomer(db *sql.DB) error {
	err := db.QueryRow("CALL PAS_SP_CREATE_CUSTOMER($1, $2, $3)", c.FirstName, c.LastName, c.CustomerPhoneNumber).Scan(&c.CustomerID)
	if err != nil {
		return err
	}

	return nil
}

// Create an order
func (o *order) createOrder(db *sql.DB) error {
	err := db.QueryRow("CALL PAS_SP_CREATE_ORDER($1, $2)", o.PizzaID, o.CustomerPhoneNumber).Scan(&o.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// Get order status
func (o *order) getStatus(db *sql.DB) error {
	return db.QueryRow("CALL PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER($1)", o.OrderID).Scan(&o.OrderStatus)
}

// Cancel an order
func (o *order) cancelOrder(db *sql.DB) error {
	return db.QueryRow("CALL PAS_SP_CANCEL_ORDER($1)", o.OrderID).Scan(&o.OrderStatus)
}

// Get orders
func (o *order) getOrders(db *sql.DB) ([]order, error) {
	// Run the query
	rows, err := db.Query(
		"SELECT customerPhoneNumber, orderId, orderTime, pizzaId, totalPrice, orderStatus FROM ORDERS WHERE customerPhoneNumber = $1 AND isDeleted = FALSE", o.CustomerPhoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a 'orders' list and append each resulting row to the list
	orders := []order{}
	for rows.Next() {
		if err := rows.Scan(&o.CustomerPhoneNumber, &o.OrderID, &o.OrderTime, &o.PizzaID, &o.TotalPrice, &o.OrderStatus); err != nil {
			return nil, err
		}
		orders = append(orders, *o)
	}

	return orders, nil
}

// Get avilable pizzas
func (p *pizza) getAvailablePizzas(db *sql.DB) ([]pizza, error) {
	// Run the query
	rows, err := db.Query(
		"SELECT pizzaId, pizzaName, pizzaPrice FROM PIZZAS WHERE isDeleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a 'pizzas' list and append each resulting row to the list
	pizzas := []pizza{}
	for rows.Next() {
		if err := rows.Scan(&p.PizzaID, &p.PizzaName, &p.PizzaPrice); err != nil {
			return nil, err
		}
		pizzas = append(pizzas, *p)
	}

	return pizzas, nil
}
