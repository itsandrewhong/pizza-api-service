# Create a new order
Notes:
* A customer can create multiple orders. An orderID is returned as the HTTP Response and this orderID will be used to track the progress of the order, cancel/update the order.
* Customers can get the list of available pizzas [HERE](doc/showPizzas.md)
*  ***Please take note of the 'orderId'. An `orederId` is required to check the order status and cancel the order.***

**URL** : `/order/add`

**Method** : `POST`

**Auth required** : Yes

**Auth constraint**
```bash
Authorization: Bearer [token]
```

**Auth example**
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs
```

**Data constraints**
```json
{
  "pizzaId": [integer], 
  "customerPhoneNumber": "[valid phone number]"
}
```

**Data example**
```json
{
  "pizzaId": 4, 
  "customerPhoneNumber": "8125984475"
}
```

## cURL Command
```bash
# Request Definition
curl -XPOST -H 'Authorization: Bearer <token>' -H "Content-type: application/json" -d '{"pizzaId": <pizzaId>, "customerPhoneNumber":"<customerPhoneNumber>"}' 'https://pizza-api-service.herokuapp.com/order/add'

# Example Request
curl -XPOST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" -d '{"pizzaId": 4, "customerPhoneNumber":"8125984475"}' 'https://pizza-api-service.herokuapp.com/order/add'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
{
  "orderId":11
}
```