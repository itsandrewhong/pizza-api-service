package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shaj13/go-guardian/auth"
	"golang.org/x/crypto/bcrypt"
)

var privKey = rsaKeySetup()

// CreateToken - Handler for creating a bearer token
func CreateToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "medium",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	// Sign the token with private key
	jwtToken, _ := token.SignedString(privKey)

	// Write to the HTTP Response
	w.Write([]byte(jwtToken))
}

// Signup -
func Signup(w http.ResponseWriter, r *http.Request) {

	// Parse and decode the request body into a new `customer` instance
	creds := &customer{}
	err := json.NewDecoder(r.Body).Decode(creds)

	// If there is something wrong with the request body, return a 400 status
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Salt and hash the password using the 'bcrypt' algorithm with salt of 10 rounds (bcrypt.DefaultCost)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

	// Insert the username and the hashed password into the DB
	// If there is any issue with inserting into the database, return a 500 error
	_, err = db.Query("INSERT INTO CUSTOMERS_DEV VALUES (DEFAULT, $1, $2, $3, $4, $5, FALSE)",
		creds.FirstName, creds.LastName, creds.CustomerPhoneNumber, creds.Username, string(hashedPassword))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ValidateUser - Verify user credential
// Handler to authenticate a user given his username and password
func ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	var err error

	if userName == "" {
		log.Println("Username cannot be empty")
		return nil, fmt.Errorf("Username cannot be empty")
	}
	if password == "" {
		log.Println("Password cannot be empty")
		return nil, fmt.Errorf("Password cannot be empty")
	}
	var retrievedPassword string

	// Retrieve the user info for the given username in the DB and
	// store the retrieved password to the 'retrievedPassword' field.
	err = db.QueryRow("select password from CUSTOMERS_DEV where username=$1", userName).Scan(&retrievedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Provided username does not exist")
			return nil, fmt.Errorf("Provided username does not exist")
		}
		log.Println("Issue with the DB")
		return nil, fmt.Errorf("Issue with the DB")
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(retrievedPassword), []byte(password)); err != nil {
		log.Println("Provided password does not match")
		return nil, fmt.Errorf("Provided password does not match")
	}

	return auth.NewDefaultUser(userName, "1", nil, nil), nil
}

// VerifyToken - Verify token
func VerifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return privKey, nil

	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := auth.NewDefaultUser(claims["sub"].(string), "", nil, nil)
		return user, nil
	}

	return nil, fmt.Errorf("Invaled token")
}
