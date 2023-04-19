# Golang Backend Starter


## Run Locally
Install dependencies

```go
go mod vendor
```

### Run Migration
```go
go run main.go migrate
```

### Start the server Locally
```go
go run main.go serve
```


## Start the server by Docker

```makefile
make development
```
or
```shell
bash run.sh
```

## Stop the server by Docker

```makefile
make clean
```

