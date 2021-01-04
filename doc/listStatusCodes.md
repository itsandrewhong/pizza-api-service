# Get the list of order status
Store employees can view the list of order status codes. This API call can be useful when an employee updates the order status.

**URL** : `/status_code/show`

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
curl -XGET -H 'Authorization: Bearer <token>' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/status_code/show'

# Example Request
curl -XGET -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs' -H "Content-type: application/json" 'https://pizza-api-service.herokuapp.com/status_code/show'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
[
  {"statusId":1,"statusName":"Order Received"},
  {"statusId":2,"statusName":"Making Your Pizza"},
  {"statusId":3,"statusName":"Ready for Pick Up"},
  {"statusId":4,"statusName":"Picked Up"},
  {"statusId":5,"statusName":"Canceled"}
]
```