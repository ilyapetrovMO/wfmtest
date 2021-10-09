# WFM Test

Available on https://wfmtest.herokuapp.com/products (free tier, expect wake up time).

Postman collection:
https://www.getpostman.com/collections/0523c21ae2b0f3c07be7

Use one of the "Login as ..." requests first to get a fresh token and populate the `TOKEN` and `USERID` variables.

-----------------------------
## Authorization scheme:
Bearer authorization scheme is used to pass a JWT token. Token carries three claims:
```
{
  "exp": 1633524512,  // expiry date in unix time
  "user_id": 1,       // user id
  "role_id": 2        // role id, where 1=manager and 2=user
}
```

example token:
`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzM4MDExODQsInVzZXJfaWQiOjEsInJvbGVfaWQiOjJ9.Xcvn3x46AyGfPojPKSYhC7Yyzai-R4X54aj00_H-oQM`

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
-----------------------------------
