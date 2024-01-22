# TODO list

- improve the readme, describe the endpoints
- add comments to functions
- potentially add an endpoint to add/delete labels directly without having the use the PUT endpoint
- add integration tests / more unit tests
- make labels a set, instead of an array on mongo, not sure if possible, otherwise solve in code

# Improvement points since last assignment

This assignment I tried to do more with middleware. The error handling and param middlewares. This way there should be less code duplication compared to last time.

# How to run in Docker

```bash
docker-compose up
```

# How to run locally

## Environment

- MONGOD_PATH: the full-path to the mongod binary on your system
- MONGO_URL: a mongo url to connect to
- USE_MEMORY_MONGO: boolean to flag if the MONGOD_PATH should be used for a memory mongo, or the MONGO_URL for a "real" mongo instance
- MAX_RETURN_ARRAY_SIZE: the max size of array returns, preventing potential memory issues
- PORT: the port where the server runs

## Testing

```bash
MONGOD_PATH=<MONGOD_PATH> go test -v ./...
```

## Running

```bash
MONGOD_PATH=<MONGOD_PATH> go run main.go
```

## Building

```bash
MONGOD_PATH=<MONGOD_PATH> go build main.go
```

# Endpoints

x DELETE /todo/:id
x GET /todo
x GET /todo/:id
x GET /todo/label/:label
x POST /todo
x PUT /todo/:id
