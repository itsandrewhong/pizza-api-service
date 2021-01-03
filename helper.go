package main

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Handle error message
func responseErrorHandler(w http.ResponseWriter, code int, message string) {
	responseWriter(w, code, map[string]string{"error": message})
}

// Write HTTP response in JSON format
func responseWriter(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Get DB connection string from file
func getConnString() string {
	connString, err := ioutil.ReadFile("cstrings.config")
	if err != nil {
		log.Fatal(err)
	}

	return string(connString)
}

// CustomerPhoneNumber cannot be empty, and must be a string consisting of ten digits
func validateCustomerPhoneNumber(i interface{}) error {
	// Compile the regex expression once
	re := regexp.MustCompile("^[0-9]{10}$")

	switch v := i.(type) {
	case order:
		return validation.ValidateStruct(&v, validation.Field(&v.CustomerPhoneNumber, validation.Required, validation.Match(re)))
	case customer:
		return validation.ValidateStruct(&v, validation.Field(&v.CustomerPhoneNumber, validation.Required, validation.Match(re)))
	}

	return errors.New("validateCustomerPhoneNumber: invalid type provided")
}

// Reads a RSA key from a file/env and decodes the key
func rsaKeySetup() []byte {
	var priv []byte
	var err error

	privString := os.Getenv("PRIVATE_KEY")

	if privString == "" {
		log.Println("No key found on cloud env, using local key")

		priv, err = ioutil.ReadFile("key/jwtRS256.key")
		if err != nil {
			log.Println("No RSA private key found, halting the application")
			panic(err)
		}
	} else {
		priv = []byte(privString)
	}

	privPem, _ := pem.Decode(priv)
	if privPem.Type != "RSA PRIVATE KEY" {
		log.Println("RSA private key is of the wrong type")
	}

	privPemBytes := privPem.Bytes
	return privPemBytes
}
