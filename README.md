```
$ git clone https://github.com/ilyapetrovMO/wfmtest.git
$ cd wfmtest
$ docker-compose up
$ docker run -v ./migrations:migrations -- netowork host migrate/migrate
    -path=/migrations/ -database postgres://postgres:postgres@localhost:8888/wfmtest up
```