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
	setupGoGuardian()

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

	//
	a.Router.HandleFunc("/v1/auth/signup", Signup).Methods("POST")

	// Route for obtaining a bearer token
	// Wrap the 'CreateToken' function with middleware to authenticate the request
	a.Router.HandleFunc("/v1/auth/token", middleware(http.HandlerFunc(CreateToken))).Methods("GET")

	// Route for accessing a protected resource that returns a book author by id
	a.Router.HandleFunc("/v1/book/{id}", middleware(http.HandlerFunc(GetBookAuthor))).Methods("GET")

	// Create a new customer
	a.Router.HandleFunc("/customer/add", a.createCustomerHandler).Methods("POST")
	// Create a new order
	a.Router.HandleFunc("/order/add", a.createOrderHandler).Methods("POST")
	// Get status of the order
	a.Router.HandleFunc("/order/show/{orderId:[0-9]+}", a.getStatusHandler).Methods("GET")
	// Cancel an order
	a.Router.HandleFunc("/order/update/{orderId:[0-9]+}", a.cancelOrderHandler).Methods("PUT")
	// Get list of orders by specific phone number
	a.Router.HandleFunc("/order/show", a.getOrdersHandler).Methods("GET")
	// Get the list of available pizzas
	a.Router.HandleFunc("/pizza/show", a.getAvailablePizzasHandler).Methods("GET")
	// Get the list of order status
	a.Router.HandleFunc("/status_code/show", a.getStatusCodeHandler).Methods("GET")
	// Update the order status
	a.Router.HandleFunc("/order/update", a.updateOrderStatusHandler).Methods("PUT")
}

// Setup Go-Guardian
func setupGoGuardian() {

	// Create an authenticator
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)

	// Dispatch the authenticator to strategies
	basicStrategy := basic.New(ValidateUser, cache)
	tokenStrategy := bearer.New(VerifyToken, cache)
	authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)

	// Cache the authentication decision to improve server performance
	authenticator.EnableStrategy(bearer.CachedStrategyKey, tokenStrategy)
}

// HTTP middleware to intercept the request and authenticate users before it reaches the final route
// This middleware will check if an access token exists and is valid. If it passes the checks, the request will proceed. If not, a 401 Authorization error is returned.

func middleware(next http.Handler) http.HandlerFunc {
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
