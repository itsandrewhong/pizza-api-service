# Update an order status 
Allows the store employees to update the order status.

**URL** : `/order/update`

**Method** : `PUT`

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
    "orderId": [valid orderId], 
    "orderStatus": [valid order status code]
}
```

**Data example**
```json
{
  "orderId": 11, 
  "orderStatus":2
}
```

## cURL Command
```bash
# Request Definition
curl -v -XPUT -H 'Authorization: Bearer <token>' -H "Content-type: application/json" -d '{"orderId": <orderId>, "orderStatus":<orderStatusCode>}' 'https://pizza-api-service.herokuapp.com/order/update'

# Example Request
curl -v -XPUT -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" -d '{"orderId": 11, "orderStatus":2}' 'https://pizza-api-service.herokuapp.com/order/update'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
{
  "orderId":"9",
  "orderStatus":"Making Your Pizza"
}
```