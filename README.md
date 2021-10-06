# WFM Test

Available on https://wfmtest.herokuapp.com/products (expect a slight pause on first request).

Postman collection:
https://www.getpostman.com/collections/bd0a787364d51580dc91

Get user token:
```
$ BODY='{"username": "user1", "password": "user1"}'
$ curl -X POST -i -d "$BODY" wfmtest.herokuapp.com/auth
```
Get manager token:
```
$ BODY='{"username": "manager1", "password": "manager1"}'
$ curl -X POST -i -d "$BODY" wfmtest.herokuapp.com/auth
```
Get all orders (manager only):
```
$ curl -H "Authorization: Bearer $MANAGER_TOK" wfmtest.herokuapp.com/orders
```

Create order (user only):
```
$ BODY='{"product_id": 1, "amount": 10}'
$ curl -X POST -H "Authorization: Bearer $USER_TOK" -d "$BODY" wfmtest.herokuapp.com/orders
```

-----------------------------
## Authorization scheme:
Bearer authorization scheme is used to pass a JWT token. Token carries three claims:
```
{
  "exp": 1633524512, // expiry date in unix time
  "UserId": 1,       // user id
  "RoleId": 2        // role id, where 1=manager and 2=user
}
```

example token:
`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzM1MjQ1MTIsIlVzZXJJZCI6MSwiUm9sZUlkIjoyfQ.Ph2Q98E9j-dlMesvCknYW-_wLQRNtv5aIjE8W_w8To4`

check its contents on www.jwt.io


------------------------------------------------------------------------
## Database structure:
https://dbdiagram.io/d/615ca2ed825b5b0146229a76

-------------------------------------------------------------------------
## Local Setup:
Requires git, docker and docker-compose.


Clone repo, build the app and start postgres:
```
$ git clone https://github.com/ilyapetrovMO/wfmtest.git
$ cd wfmtest
$ docker-compose up -d
```
 Populate the DB:
```
$ docker run -v $PWD/migrations:/migrations --network host migrate/migrate
    -path=/migrations/ -database 'postgres://postgres:postgres@localhost:8888/wfmtest?sslmode=disable' up
```

If `docker-compose logs` reports that server_1 had an unexpected error, restart with:
```
docker-compose down
docker-compose up
``` 