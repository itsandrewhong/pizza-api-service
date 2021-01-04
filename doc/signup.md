# Add - Create a new customer account 

Notes:
* A phone number is unique. A customer can only create one account using the same phone number. 
* A phone number must be a non-null string consisting exactly ten digits without country code (e.g. 8125984475).

**URL** : `/customer/add`

**Method** : `POST`

**Auth required** : NO

**Data constraints**
```json
{
  "firstName": "[unicode 50 chars max]", 
  "lastName": "[unicode 50 chars max]", 
  "customerPhoneNumber": "[numeric value 10 integers max]", 
  "username": "[unicode 62 chars max]", 
  "password": "[unicode 60 chars max]",
}
```

```json
{
  "firstName": "[firstName in plain text]", 
  "lastName": "[lastName in plain text]", 
  "customerPhoneNumber": "[valid phone number]", 
  "username": "[username in plain text]", 
  "password": "[password in plain text]"
}
```

**Data example**
```json
{
  "firstName": "Carl", 
  "lastName": "Raymond", 
  "customerPhoneNumber": "8125984475", 
  "username": "carl@gmail.com", 
  "password": "mysecurepassword"
}
```

## cURL Command
```bash
# Request Definition
curl -v -XPOST -H "Content-type: application/json" 
-d '{"firstName":"<firstName>", "lastName":"<lastName>", "customerPhoneNumber":"<customerPhoneNumber>", "username": "<username>", "password": "<password>"}' 'https://pizza-api-service.herokuapp.com/customer/add'

# Example Request
curl -v -XPOST -H "Content-type: application/json" -d '{"firstName":"Carl", "lastName":"Raymond", "customerPhoneNumber":"8125984475", "username": "carl@gmail.com", "password": "mysecurepassword"}' 'https://pizza-api-service.herokuapp.com/customer/add'
```

## Success Response
**Code** : `200 OK`

**Content example**

```json
{
    "customerPhoneNumber": "8125984475"
}
```

