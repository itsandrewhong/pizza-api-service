# Token - Obtain a user access token
The -u option in cURL performs Basic Authentication, where you can effectively "login" to an API by using a username and password.

**URL** : `/auth/token`

**Method** : `GET`

**Auth required** : Yes

**Auth constraints**
```bash
  "[unicode 62 chars max]":"[unicode 60 chars max]"
```

**Auth example**
```bash
"alice@gmail.com:mysecurepassword"
```

## cURL Command
```bash
# Request Definition
curl -XGET 'https://pizza-api-service.herokuapp.com/auth/token' -u "<username>:<password>"

# Example Request
curl -XGET 'https://pizza-api-service.herokuapp.com/auth/token' -u "alice@gmail.com:mysecurepassword"
```

## Success Response
**Code** : `200 OK`

**Content example**

```bash
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2MDk3ODgwMjIsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.hrLAEbHKLHrTXG7_9TVot8Dubq2hHia5khMQeTUqJLs
```

