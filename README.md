# pizza-api-service
A simple REST API application for Pizza Ordering System using `Golang`, `PostgreSQL` as database, and `Gorilla Mux` for routing

## Pizza Ordering System business logic:
- [MVP] **Create** a new customer in response to a valid `POST` request at `/customer/add` with first name, last name, and customer phone number.
- [MVP] **Create** a new order in response to a valid `POST` request at `/order/add` with pizzaId, and customer phone number.
- [MVP] **Fetch** the order status in response to a valid `GET` request at `/order/show/` with orderId.
- **Fetch** the list of available pizzas in response to a valid `GET` request at `/pizza/show`.
- **Cancel** an order in response to a valid `PUT` request at `/order/update/` with orderId.
- **Fetch** the list of order in response to a valid `GET` request at `/order/show/` with customer phone number.

## App Dependencies
1. 'mux' - Gorilla Mux router, used to create complex routing and managing requests
2. 'pq' - PostgreSQL driver, used to store the data
3. 'ozzo-validation' - Input validation, used to validate the user input (phone number, name, etc.)

## File Structure
* `main.go`: Initializes DB connection and Runs the application.
* `app.go`: Contains the API business logic, definition to connect app with the DB, and definition to run the application.
* `model.go`: Setup structs to connect Golang with DB(Postgres) and interacts with the Database.

## Test the application via cURL commands
- ***NOTE: The application is hosted on Heroku (free-tier). The DB will sleep after a half hour of inactivity, and it causes a delay of a few seconds for the first request upon waking.***
```bash
# Create a customer
# Note: A phone number must be a non-null string consisting exactly ten digits without country code (e.g. +1)
curl -v -XPOST -H "Content-type: application/json" -d '{"firstName":"Carl", "lastName":"Raymond", "customerPhoneNumber":"8481259874"}' 'https://pizza-api-service.herokuapp.com/customer/add'

# Get the list of available pizzas
curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/pizza/show'

# Create an order
curl -v -XPOST -H "Content-type: application/json" -d '{"pizzaId": 2, "customerPhoneNumber":"8481259874"}' 'https://pizza-api-service.herokuapp.com/order/add'

# Check status of the order
# curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/<orderId>'
curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/2'

# Cancel an order
# curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/<orderId>'
curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/2'

# Get order info
curl -v -XGET -H "Content-type: application/json" -d '{"customerPhoneNumber":"8481259874"}' 'https://pizza-api-service.herokuapp.com/order/show'
```

## Database: Stored Procedure Definitions
```sql
-- Create a customer (PAS_SP_CREATE_CUSTOMER)
CREATE PROCEDURE PAS_SP_CREATE_CUSTOMER(
	IN p_firstName VARCHAR(50),
	IN p_lastName VARCHAR(50),
	IN p_customerPhoneNumber VARCHAR(20),
	INOUT _customerId INTEGER DEFAULT null
)
LANGUAGE SQL
AS $$
	INSERT INTO CUSTOMERS VALUES (DEFAULT, TRIM(p_firstName), TRIM(p_lastName), TRIM(p_customerPhoneNumber), FALSE) RETURNING customerId;
$$;

-- Create an order (PAS_SP_CREATE_ORDER)
CREATE PROCEDURE PAS_SP_CREATE_ORDER(
	IN p_pizzaId INTEGER,
	IN p_customerPhoneNumber VARCHAR(20),
	INOUT _orderId INTEGER DEFAULT null
)
LANGUAGE SQL
AS $$
	INSERT INTO ORDERS VALUES (DEFAULT, p_pizzaId, CURRENT_TIMESTAMP, TRIM(p_customerPhoneNumber), 'Order Received', 
  ROUND(((SELECT p.pizzaPrice FROM PIZZAS p where p.pizzaId = p_pizzaId) * 1.0625), 2), 
  FALSE) RETURNING orderId;
$$;

-- Fetch an order status (PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER)
CREATE PROCEDURE PAS_SP_GET_ORDER_STATUS_BY_ORDERNUMBER(
	IN p_orderId INTEGER,
	INOUT _orderStatus VARCHAR(30) DEFAULT null
)
LANGUAGE SQL
AS $$
	SELECT o.orderstatus FROM ORDERS o WHERE o.orderId = p_orderId AND o.isDeleted = FALSE;
$$;

-- Cancel an order (PAS_SP_CANCEL_ORDER)
CREATE PROCEDURE PAS_SP_CANCEL_ORDER(
	IN p_orderId INTEGER,
	INOUT _orderStatus VARCHAR(30) DEFAULT null
)
LANGUAGE SQL
AS $$
	UPDATE ORDERS SET orderStatus='Canceled' WHERE orderId = p_orderId;
	SELECT o.orderstatus FROM ORDERS o WHERE o.orderId = p_orderId;
$$;
```