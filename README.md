# pizza-api-service
A simple REST API application for Pizza Ordering System using `Golang`, `PostgreSQL`, and `Gorilla Mux`.

## Pizza Ordering System business logic:
- [MVP] **Create** a new customer in response to a valid `POST` request at `/customer/add` with first name, last name, and customer phone number.
- [MVP] **Create** a new order in response to a valid `POST` request at `/order/add` with pizzaId, and customer phone number.
- [MVP] **Fetch** the order status in response to a valid `GET` request at `/order/show/<orderId>` with order ID.
- **Fetch** the list of available pizzas in response to a valid `GET` request at `/pizza/show`.
- **Cancel** an order in response to a valid `PUT` request at `/order/update/<orderId>` with order ID.
- **Fetch** the list of orders by specific phone number in response to a valid `GET` request at `/order/show` with customer phone number.
- **Fetch** the list of order status code in response to a valid `GET` request at `/status_code/show`. (Store Only)
- **Update** an order in response to a valid `PUT` request at `/order/update` with order ID and order status code. (Store Only)

## App Dependencies
1. `mux` - Gorilla Mux router, used to create complex routing and managing requests
2. `pq` - PostgreSQL driver, used to store the data
3. `ozzo-validation` - Input validation, used to validate the user input (phone number, name, etc.)

## File Structure
* `main.go`: Initializes DB connection and Runs the application.
* `app.go`: Contains the API business logic, definition to connect app with the DB, and definition to run the application.
* `model.go`: Setup structs to connect Golang with DB(Postgres) and interacts with the Database.

# Example) Test the application via cURL commands
- ***NOTE: The application is hosted on Heroku (free-tier). The DB will sleep after a half hour of inactivity, and it causes a delay of a few seconds for the first request upon waking.***

## Create a new customer
* A phone number is unique. A customer can only create one account using the same phone number. 
* A phone number must be a non-null string consisting exactly ten digits without country code (e.g. 8125984475).
```bash
# Request
curl -v -XPOST -H "Content-type: application/json" -d '{"firstName":"Carl", "lastName":"Raymond", "customerPhoneNumber":"8125984475"}' 'https://pizza-api-service.herokuapp.com/customer/add'

# Response
{"customerPhoneNumber":"8125984475"}*
```

## Get the list of available pizzas
* This application uses a numeric pizzaId to create an order. This way, the store can update the pizza info when needed.
* Displays the list of available pizzas.
```bash
# Request
curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/pizza/show'

# Response
[{"pizzaId":1,"pizzaName":"Cheese Pizza","pizzaPrice":6.99},
{"pizzaId":2,"pizzaName":"Veggie Pizza","pizzaPrice":7.99},
{"pizzaId":3,"pizzaName":"Pepperoni Pizza","pizzaPrice":6.99},
{"pizzaId":4,"pizzaName":"Meat Pizza","pizzaPrice":7.99},
{"pizzaId":5,"pizzaName":"Margherita Pizza","pizzaPrice":8.99},
{"pizzaId":6,"pizzaName":"BBQ Chicken Pizza","pizzaPrice":8.99},
{"pizzaId":7,"pizzaName":"Hawaiian Pizza","pizzaPrice":7.99},
{"pizzaId":8,"pizzaName":"Buffalo Pizza","pizzaPrice":10.99},
{"pizzaId":9,"pizzaName":"Supreme Pizza","pizzaPrice":12.99}]
```

## Create a new order
* A customer can create multiple orders.
* An orderID is returned as the HTTP Response and this orderID will be used to track the progress of the order, cancel/update the order.
```bash
# Request
curl -v -XPOST -H "Content-type: application/json" -d '{"pizzaId": 4, "customerPhoneNumber":"8125984475"}' 'https://pizza-api-service.herokuapp.com/order/add'

# Response
{"orderId":9}
```

## Check status of the order
* Allows to check the status of the order given the orderId in the URL
```bash
# Request
# curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/<orderId>'
curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/9'

# Response
{"orderStatus":"Order Received"}
```

## Cancel an order
* A cusotmer may have changered his/her mind, the application allows to cancel an order given the orderId in the URL
```bash
# Request
# curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/<orderId>'
curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/9'

# Response
{"orderStatus":"Canceled"}
```

# Get the list of orders by specific phone number
* Customers can view their order history with orderId, order status, etc.
```bash
# Request
curl -v -XGET -H "Content-type: application/json" -d '{"customerPhoneNumber":"8125984475"}' 'https://pizza-api-service.herokuapp.com/order/show'

# Response
[{"orderId":9,"pizzaId":4,"orderTime":"2020-12-27T21:56:41.636116Z","customerPhoneNumber":"8125984475","orderStatus":"Canceled","totalPrice":8.49}]
```

# Store Only: Get the list of order status
* Displays the list of order status code
```bash
# Request
curl -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/status_code/show'

# Response
[{"statusId":1,"statusName":"Order Received"},
{"statusId":2,"statusName":"Making Your Pizza"},
{"statusId":3,"statusName":"Ready for Pick Up"},
{"statusId":4,"statusName":"Picked Up"},
{"statusId":5,"statusName":"Canceled"}]
```

# Store Only: Update the order status 
* Allows the store employees to update the order status.
```bash
# Request
# curl -v -XPUT -H "Content-type: application/json" -d '{"orderId": <orderId>, "orderStatus":<orderStatusCode>}' 'https://pizza-api-service.herokuapp.com/order/update'

curl -v -XPUT -H "Content-type: application/json" -d '{"orderId": 9, "orderStatus":2}' 'https://pizza-api-service.herokuapp.com/order/update'

# Response
{"orderId":"9","orderStatus":"Making Your Pizza"}
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
```