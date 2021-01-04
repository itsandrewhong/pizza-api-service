# Show available pizzas

Notes:
* This API call takes in the token as a header. Token can be created [HERE](token.md)
* This application uses a numeric pizzaId to create an order. This way, the store can update the pizza info when needed.
* Displays the list of available pizzas.

**URL** : `/pizza/show`

**Method** : `GET`

**Auth required** : Yes

**Auth constraint**
```bash
Authorization: Bearer [token]
```

**Auth example**
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs
```

## cURL Command
```bash
# Request Definition
curl -XGET -H 'Authorization: Bearer <token>' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/pizza/show'

# Example Request
curl -XGET -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/pizza/show'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
[
  {"pizzaId":1, "pizzaName":"Cheese Pizza", "pizzaPrice":6.99},
  {"pizzaId":2, "pizzaName":"Veggie Pizza", "pizzaPrice":7.99},
  {"pizzaId":3, "pizzaName":"Pepperoni Pizza", "pizzaPrice":6.99},
  {"pizzaId":4, "pizzaName":"Meat Pizza", "pizzaPrice":7.99},
  {"pizzaId":5, "pizzaName":"Margherita Pizza", "pizzaPrice":8.99},
  {"pizzaId":6, "pizzaName":"BBQ Chicken Pizza", "pizzaPrice":8.99},
  {"pizzaId":7, "pizzaName":"Hawaiian Pizza", "pizzaPrice":7.99},
  {"pizzaId":8, "pizzaName":"Buffalo Pizza", "pizzaPrice":10.99},
  {"pizzaId":9, "pizzaName":"Supreme Pizza", "pizzaPrice":12.99}
]
```