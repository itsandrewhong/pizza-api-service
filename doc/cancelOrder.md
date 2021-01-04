# Cancel an order
A cusotmer may have changered his/her mind, the application allows to cancel an order given the orderId in the URL

**URL** : `/order/update/{orderId:[0-9]+}`

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

## cURL Command
```bash
# Request Definition
curl -v -XPUT -H 'Authorization: Bearer <token>' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/<orderId>'

# Example Request
curl -v -XPUT -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/order/update/11'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
{
  "orderStatus" : "Canceled"
}
```