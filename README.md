# pizza-api-service
A simple REST API application for Pizza Ordering System with token-based authentication.

## Pizza Ordering System business logic:
- [MVP] **Create** a new customer in response to a valid `POST` request at `/customer/add` with first name, last name, customer phone number, username, and password.
- [MVP] **Create** a new order in response to a valid `POST` request at `/order/add` with pizzaId, and customer phone number.
- [MVP] **Fetch** the order status in response to a valid `GET` request at `/order/show/<orderId>` with order ID.
- **Fetch** the list of available pizzas in response to a valid `GET` request at `/pizza/show`.
- **Cancel** an order in response to a valid `PUT` request at `/order/update/<orderId>` with order ID.
- **Fetch** the list of orders by specific phone number in response to a valid `GET` request at `/order/show` with customer phone number.
- **Fetch** the list of order status code in response to a valid `GET` request at `/status_code/show`. (Store Only)
- **Update** an order in response to a valid `PUT` request at `/order/update` with order ID and order status code. (Store Only)

## User Authentication
This application uses:
1. `bcrypt` algorithm to hash and salt the customer's password
2. `jwt-go` to implement a stateless authentication; create a token, sign it with the server's secret key (token is valid for 24 hours), and validate/verify the token,
3. `go-guardian` to authenticate requests and cache the authentication decisions

## App Dependencies
1. `mux` - Gorilla Mux router, used to create complex routing and managing requests
2. `pq` - PostgreSQL driver, used to store the data
3. `ozzo-validation` - Input validation, used to validate the user input (phone number, name, etc.)
4. `go-guardian` - A simple API authentication
5. `jwt-go` - Used to create and verify the token
6. `bcrypt` - Used to encrypt (hash and salt) the user password

## File Structure
* `main.go`: Initializes DB connection and Runs the application.
* `app.go`: Contains routes, definition to connect app with the DB, and definition to run the application.
* `handler.go`: Contains the API business logic.
* `model.go`: Setup structs to connect Golang with DB(Postgres) and interacts with the Database.
* `authHandler.go`: Contains the functions to create, validate, and verify a token.
* `helper.go`: Contains the helper functions that support the application.


# API Calls (Test API via cURL commands)
NOTE:
* ***The application is hosted on Heroku (free-tier). The DB will sleep after a half hour of inactivity, and it causes a delay of a few seconds for the first request upon waking.***
* ***After a customer account has been created, a customer must [Obtain user access token](token.md) in order to view list of pizzas and make order related calls.***

## General Process Flow
1. [Create customer](doc/signup.md)
2. [Obtain user access token](doc/token.md)
3. [Show available pizzas](doc/showPizzas.md)
4. [Create new order](doc/createOrder.md)
5. ...

## Open Endpoints
Open endpoints require no Authentication.
* [Create customer](doc/signup.md) : `POST /customer/add`

## Endpoints that require Authentication
Closed endpoints require a basic authentication with username and password, or a valid Token to be included in the header of the request. A Token can be acquired from the `Create customer` view above.

### Current User related
Endpoint for obtaining the token that the Authenticated User has permissions to access.
* [Obtain user access token](doc/token.md) : `GET /auth/token`

### Pizza related
Endpoint for viewing the Pizzas that the Authenticated User has permissions to access.
* [Show available pizzas](doc/showPizzas.md) : `GET /pizza/show`

### Order related
Endpoints for viewing and manipulating the Orders that the Authenticated User has permissions to access.
* [Create a new order](doc/createOrder.md) : `POST /order/add`
* [Check status of the order](doc/getOrderStatus.md) : `GET /order/show/{orderId:[0-9]+}`
* [Cancel an order](doc/cancelOrder.md) : `PUT /order/update/{orderId:[0-9]+}`
* [Show orders by specific phone number](doc/getOrdersByPhoneNumber.md) : `GET /order/show`

***Below calls are made ONLY by store employees. Current version allows customers to call below API calls for testing purposes and simplicity but will be only applicable to store employees in the later versions.***
* [Show list of order status](doc/listStatusCodes.md) : `GET /status_code/show`
* [Update order status ](doc/updateOrderStatus.md) : `PUT /order/update`


# Database: Stored Procedure Definitions
```sql
-- Create a customer (PAS_SP_CREATE_CUSTOMER)
CREATE PROCEDURE PAS_SP_CREATE_CUSTOMER(
	IN p_firstName VARCHAR(50),
	IN p_lastName VARCHAR(50),
	IN p_customerPhoneNumber VARCHAR(20),
	IN p_userName VARCHAR(62),
	IN p_password VARCHAR(60),
	INOUT _customerId INTEGER DEFAULT null
)
LANGUAGE SQL
AS $$
	INSERT INTO CUSTOMERS VALUES (DEFAULT, TRIM(p_firstName), TRIM(p_lastName), TRIM(p_customerPhoneNumber), TRIM(p_userName), TRIM(p_password), FALSE) RETURNING customerId;
$$;

-- Create an order (PAS_SP_CREATE_ORDER)
CREATE PROCEDURE PAS_SP_CREATE_ORDER(
	IN p_pizzaId INTEGER,
	IN p_customerPhoneNumber VARCHAR(20),
	INOUT _orderId INTEGER DEFAULT null
)
LANGUAGE SQL
AS $$
	INSERT INTO ORDERS VALUES (DEFAULT, p_pizzaId, CURRENT_TIMESTAMP, TRIM(p_customerPhoneNumber), 1, 
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
	SELECT sc.statusName FROM ORDERS AS o INNER JOIN ORDER_STATUS_CODES AS sc ON o.statusId = sc.statusId WHERE o.orderId = p_orderId;
$$;

-- Cancel an order (PAS_SP_CANCEL_ORDER)
CREATE PROCEDURE PAS_SP_CANCEL_ORDER(
	IN p_orderId INTEGER,
	INOUT _orderStatus VARCHAR(30) DEFAULT null
)
LANGUAGE SQL
AS $$
	UPDATE ORDERS SET statusId = 5 WHERE orderId = p_orderId;	
	SELECT sc.statusName FROM ORDERS AS o INNER JOIN ORDER_STATUS_CODES AS sc ON o.statusId = sc.statusId WHERE o.orderId = p_orderId;
$$;

-- Update an order status (PAS_SP_UPDATE_ORDER_STATUS)
CREATE PROCEDURE PAS_SP_UPDATE_ORDER_STATUS(
	IN p_orderId INTEGER,
	IN p_statusId INTEGER,
	INOUT _orderStatus VARCHAR(30) DEFAULT null
)
LANGUAGE SQL
AS $$
	UPDATE ORDERS SET statusId = p_statusId WHERE orderId = p_orderId;
	SELECT sc.statusName FROM ORDERS AS o INNER JOIN ORDER_STATUS_CODES AS sc ON o.statusId = sc.statusId WHERE o.orderId = p_orderId;
$$;

-- Fetch a password given the username (PAS_SP_GET_CUSTOMER_PASSWORD)
CREATE PROCEDURE PAS_SP_GET_CUSTOMER_PASSWORD(
	IN p_userName VARCHAR(62),
	INOUT _password VARCHAR(60) DEFAULT null
)
LANGUAGE SQL
AS $$
	SELECT password FROM CUSTOMERS WHERE username = p_userName;
$$;
```