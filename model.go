package main

import (
	"database/sql"
	"errors"
)

type customer struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type order struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *customer) createCustomer(db *sql.DB) error {
	return errors.New("Not Implemented Yet")
}

func (p *order) createOrder(db *sql.DB) error {
	return errors.New("Not Implemented Yet")
}

func (p *order) getStatus(db *sql.DB) error {
	return errors.New("Not Implemented Yet")
}
