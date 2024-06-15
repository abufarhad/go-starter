Welcome to My Clean Architecture Boilerplate! 🚀

The goal is to provide you with a solid and scalable foundation for developing applications that adhere to the principles of clean architecture. With a focus on modularity, flexibility, and maintainability, our boilerplate includes essential features like Dependency Injection, Database Integration, Logging, and Error Handling.

Key Features:

- Clean Architecture Principles
- Modularity and Scalability
- Dependency Injection
- Database Integration
- Logging
- Error Handling

Join this exciting journey of building robust and high-quality software solutions. Let's create innovative applications that make a difference in the world! 🌟


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

## Swagger Docs

```go
http://localhost:8080/api/v1/docs/index.html
```

## Version
`v2.0.0 => env`   
`v1.0.0 => consul`
