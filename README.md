# WFM Test

## Setup:
Requires git, docker and docker-compose.
### Clone repo, build the app and start postgres:
```
$ git clone https://github.com/ilyapetrovMO/wfmtest.git
$ cd wfmtest
$ docker-compose up -d
```
### Populate the DB:
```
$ docker run -v $PWD/migrations:/migrations --network host migrate/migrate
    -path=/migrations/ -database 'postgres://postgres:postgres@localhost:8888/wfmtest?sslmode=disable' up
```

If `docker-compose logs` reports that server_1 had an unexpected error, restart with:
```
docker-compose down
docker-compose up
``` 