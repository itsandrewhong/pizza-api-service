# pizza-api-service
A simple REST API application for Pizza Ordering System using `Golang`, `PostgreSQL` as database, and `Gorilla Mux` for routing

# Test the application via cURL commands
- ***NOTE: The application is hosted on Heroku (free-tier). The DB will sleep after a half hour of inactivity, and it causes a delay of a few seconds for the first request upon waking.***
```bash
# Create a customer


# Create an order


# Check status of the order

```

## Pizza Ordering System business logic:
- [MVP] Create a new custoer with first name, last name, and customer phone number.
- [MVP] Create a new order with pizzaId, and customer phone number.
- [MVP] Fetch the order status with orderId.

## File Structure
`main.go`: Initializes DB connection and Runs the application.
`app.go`: Contains the API business logic, definition to connect app with the DB, and definition to run the application.
`model.go`: Setup structs to connect Golang with DB(Postgres) and interacts with the Database.

## App Dependencies
1. 'mux' - Gorilla Mux router
2. 'pq' - PostgreSQL driver

## Database: Stored Procedure Definitions
```sql


```