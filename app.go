package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize a DB connection and initialize the router
func (a *App) Initialize() {
	a.initDB()

	// Init GoGuardian
	a.setupGoGuardian()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run the application
func (a *App) Run(addr string) {
	log.Println("listen on", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Get the connection string from Heroku
// If above operation fails, set connection string manually for local
func (a *App) initDB() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = getConnString()
		log.Println("Using local connection")
	}

	var err error
	a.DB, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize routes
func (a *App) initializeRoutes() {
	// Route for creating a new customer
	a.Router.HandleFunc("/customer/add", a.createCustomerHandler).Methods("POST")

	// Route for obtaining a bearer token given the username and password
	a.Router.HandleFunc("/auth/token", middleware(createTokenHandler))
	http.Handle("/auth/token", a.Router)

	// *** Wrap below functions with middleware to authenticate the request before they reach to the final route ***

	// Route for creating a new order
	a.Router.HandleFunc("/order/add", middleware(a.createOrderHandler)).Methods("POST")
	http.Handle("/order/add", a.Router)

	// Route for retrieving status of the order
	a.Router.HandleFunc("/order/show/{orderId:[0-9]+}", middleware(a.getStatusHandler)).Methods("GET")
	http.Handle("/order/show/{orderId:[0-9]+}", a.Router)

	// Route for canceling an order
	a.Router.HandleFunc("/order/update/{orderId:[0-9]+}", middleware(a.cancelOrderHandler)).Methods("PUT")
	http.Handle("/order/update/{orderId:[0-9]+}", a.Router)

	// Route for retrieving list of orders by specific phone number
	a.Router.HandleFunc("/order/show", middleware(a.getOrdersHandler)).Methods("GET")
	http.Handle("/order/show", a.Router)

	// Route for retrieving the list of available pizzas
	a.Router.HandleFunc("/pizza/show", middleware(a.getAvailablePizzasHandler)).Methods("GET")
	http.Handle("/pizza/show", a.Router)

	// Route for retrieving the list of order status
	a.Router.HandleFunc("/status_code/show", middleware(a.getStatusCodeHandler)).Methods("GET")
	http.Handle("/status_code/show", a.Router)

	// Route for updating the order status
	a.Router.HandleFunc("/order/update", middleware(a.updateOrderStatusHandler)).Methods("PUT")
	http.Handle("/order/update", a.Router)
}

// Setup Go-Guardian
func (a *App) setupGoGuardian() {
	// Create an authenticator
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)

	// Dispatch the authenticator to strategies
	basicStrategy := basic.New(a.ValidateUserHandler, cache)
	tokenStrategy := bearer.New(VerifyTokenHandler, cache)
	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)

	// Cache the authentication decision to improve server performance
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

// HTTP middleware to intercept the request and authenticate users before it reaches the final route
// Checks if an access token exista and is valid. If it passes the checks, the request will proceed.
func middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}
