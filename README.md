# shorturl-golang

A simple but scalable ShortURL service written in Golang.

## Build Executables
```sh
go build cmd/server/server.go
go build cmd/cleaner/cleaner.go
```

## Configuration through Environment Variables
Configurations are loaded from the environment variables, if not provided, then the default values will be used.
* `LISTEN_ADDR`: The address that the server listens on, defaults to `:8000`
* `BASE_URL`: The base URL of the service, defaults to `http://localhost:8000`
* `MAX_TOKEN`: The maximum amount of tokens to be generated and stored offline, defaults to `5000`
* `MAX_URL_LEN`: The length limit of url string, defaults to `1024`
* `MONGODB_URI`: The uri to connect to mongodb, defaults to `mongodb://localhost:27017`
* `DATABASE`: The database to use, defaults to `shorturl`
* `RECORD_COLLECTION`: The collection to store records, defaults to `records`