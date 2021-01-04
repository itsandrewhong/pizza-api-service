# Get the list of orders by specific phone number
Customers can view their order history with orderId, order status, etc.

**URL** : `/order/show`

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

**Data constraints**
```json
{
  "customerPhoneNumber": "[valid phone number]"
}
```

**Data example**
```json
{
  "customerPhoneNumber": "8125984475"
}
```

## cURL Command
```bash
# Request Definition
curl -v -XGET -H 'Authorization: Bearer <token>' -H "Content-type: application/json" -d '{"customerPhoneNumber":"<customerPhoneNumber>"}' 'https://pizza-api-service.herokuapp.com/order/show'

# Example Request
curl -v -XGET -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" -d '{"customerPhoneNumber":"8125984475"}' 'https://pizza-api-service.herokuapp.com/order/show'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
[
  {
    "orderId":11,
    "pizzaId":4,
    "orderTime":"2020-12-27T21:56:41.636116Z",
    "customerPhoneNumber":"8125984475",
    "orderStatus":"Canceled",
    "totalPrice":8.49
  }
]
```