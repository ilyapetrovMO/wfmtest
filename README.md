## Clone repo, build the app and start postgres:
```
$ git clone https://github.com/ilyapetrovMO/wfmtest.git
$ cd wfmtest
$ docker-compose up -d
```
## Populate the DB:
```
$ docker run -v $PWD/migrations:/migrations --network host migrate/migrate
    -path=/migrations/ -database 'postgres://postgres:postgres@localhost:8888/wfmtest?sslmode=disable' up
```