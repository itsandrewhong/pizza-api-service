# Check status of the order
Allows to check the status of the order given the orderId in the URL

**URL** : `/order/show/{orderId:[0-9]+}`

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
curl -v -XGET -H 'Authorization: Bearer <token>' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/<orderId>'

# Example Request
curl -v -XGET -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/show/11'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
{
  "orderStatus" : "Order Received"
}
```