package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shaj13/go-guardian/auth"

	"golang.org/x/crypto/bcrypt"
)

// Setup global RSA private key
var privKey = rsaKeySetup()

// CreateTokenHandler - Handler for creating a bearer token
// Token is valid for 24 hours.
func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-app",
		"sub": "medium",
		"aud": "any",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with private key
	jwtToken, _ := token.SignedString(privKey)

	// Write to the HTTP Response
	w.Write([]byte(jwtToken))
}

// ValidateUserHandler - Validate user credential
// Handler to authenticate a user given his username and password
func (a *App) ValidateUserHandler(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	var err error

	if userName == "" {
		log.Println("Username cannot be empty")
		return nil, fmt.Errorf("Username cannot be empty")
	}
	if password == "" {
		log.Println("Password cannot be empty")
		return nil, fmt.Errorf("Password cannot be empty")
	}

	// Retrieves the hashed password from the DB given the username
	var c customer
	c.Username = userName

	// Fetch hashed password from the DB
	if err := c.getCustomerPassword(a.DB); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Provided username does not exist")
			return nil, fmt.Errorf("Provided username does not exist")
		}
		log.Println(err)
		log.Println("Issue with the DB")
		return nil, fmt.Errorf("Issue with the DB")
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
		log.Println("Provided password does not match")
		return nil, fmt.Errorf("Provided password does not match")
	}

	return auth.NewDefaultUser(userName, "1", nil, nil), nil
}

// VerifyTokenHandler - Handler to verify the given token
func VerifyTokenHandler(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {

	// Parse the token string and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return privKey, nil

	})

	if err != nil {
		return nil, err
	}

	// Authenticate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := auth.NewDefaultUser(claims["sub"].(string), "", nil, nil)
		return user, nil
	}

	return nil, fmt.Errorf("Invaled token")
}
