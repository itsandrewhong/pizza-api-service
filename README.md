# pizza-api-service
A simple REST API application for Pizza Ordering System using `Golang`, `PostgreSQL` as database, and `Gorilla Mux` for routing

## Pizza Ordering System business logic:
- [MVP] **Create** a new customer in response to a valid `POST` request at `/customer/add` with first name, last name, and customer phone number.
- [MVP] **Create** a new order in response to a valid `POST` request at `/order/add` with pizzaId, and customer phone number.
- [MVP] **Fetch** the order status in response to a valid `GET` request at `/order/show/<orderId>` with order ID.
- **Fetch** the list of available pizzas in response to a valid `GET` request at `/pizza/show`.
- **Cancel** an order in response to a valid `PUT` request at `/order/update/<orderId>` with order ID.
- **Fetch** the list of orders by specific phone number in response to a valid `GET` request at `/order/show` with customer phone number.
- **Fetch** the list of order status code in response to a valid `GET` request at `/status_code/show` with customer phone number.
- **Update** an order in response to a valid `PUT` request at `/order/update` with order ID and order status code.

## App Dependencies
1. 'mux' - Gorilla Mux router, used to create complex routing and managing requests
2. 'pq' - PostgreSQL driver, used to store the data
3. 'ozzo-validation' - Input validation, used to validate the user input (phone number, name, etc.)

## File Structure
* `main.go`: Initializes DB connection and Runs the application.
* `app.go`: Contains the API business logic, definition to connect app with the DB, and definition to run the application.
* `model.go`: Setup structs to connect Golang with DB(Postgres) and interacts with the Database.

# Test the application via cURL commands
- ***NOTE: The application is hosted on Heroku (free-tier). The DB will sleep after a half hour of inactivity, and it causes a delay of a few seconds for the first request upon waking.***

## Create a new customer
```bash
# Request
# Note: A phone number must be a non-null string consisting exactly ten digits without country code (e.g. +1)
curl -v -XPOST -H "Content-type: application/json" -d '{"firstName":"Carl", "lastName":"Raymond", "customerPhoneNumber":"8485941259"}' 'https://pizza-api-service.herokuapp.com/customer/add'

# Response
    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 2
```


## Get the list of available pizzas
* Displays the list of available pizzas
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
```bash
# Request
curl -v -XPOST -H "Content-type: application/json" -d '{"pizzaId": 2, "customerPhoneNumber":"8485941259"}' 'https://pizza-api-service.herokuapp.com/order/add'

# Response
```


## Check status of the order
* Allows to check the status of the order
```bash
# Request
# curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/<orderId>'
curl -v -XGET -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/2'

# Response
```

## Cancel an order
* Allows to cancel an order
```bash
# Request
# curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/<orderId>'
curl -v -XPUT -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/2'

# Response
```

# Get the list of orders by specific phone number
* Displays the list of orders by a specific phone number
```bash
# Request
curl -v -XGET -H "Content-type: application/json" -d '{"customerPhoneNumber":"8485941259"}' 'https://pizza-api-service.herokuapp.com/order/show'


# Response
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

curl -v -XPUT -H "Content-type: application/json" -d '{"orderId": 2, "orderStatus":5}' 'https://pizza-api-service.herokuapp.com/order/update'

# Response
{"orderId":"2","orderStatus":"Canceled"}
```


## Database: Stored Procedure Definitions
```sql


```