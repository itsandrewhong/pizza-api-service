package main

import (
	"database/sql"
)

type customer struct {
	CustomerID          int    `json:"customerId"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	IsDeleted           bool   `json:"isDeleted"`
}

type order struct {
	OrderID             int     `json:"orderId"`
	PizzaID             int     `json:"pizzaId"`
	OrderTime           string  `json:"orderTime"`
	CustomerPhoneNumber string  `json:"customerPhoneNumber"`
	OrderStatus         string  `json:"orderStatus"`
	TotalPrice          float64 `json:"totalPrice"`
	IsDeleted           bool    `json:"isDeleted"`
}

// Create a customer
func (c *customer) createCustomer(db *sql.DB) error {
	// err := db.QueryRow("INSERT INTO CUSTOMERS(firstName, lastName, customerPhoneNumber, isDeleted) VALUES ($1, $2, $3, FALSE) RETURNING customerId", c.FirstName, c.LastName, c.CustomerPhoneNumber).Scan(&c.CustomerID)

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
