# Golang, postgres and redis example


## Setup
run all server
```
docker-compose up
```

create table
```
psql -h localhost -U root -d test < structure.sql
```

## Run app
```
go run main.go
```

## Benchmark
```
go test -benchmem -bench . app
```