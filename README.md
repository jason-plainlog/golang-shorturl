# golang-shorturl

A simple but scalable ShortURL service written in Golang.

![](images/arch.jpg)

## Services
* `server`
    - API backend server
    - Includs a token collection service that collects unused ids before getting request
    - May have multiple server running at the same time
* `cleaner`
    - Small service that periodically clean up the expired records from database.

## Quickstart
Run `docker compose up` or `docker-compose up`, and voila! The service will be served on localhost:8000.

## Build Executables
```sh
go build cmd/server/server.go
go build cmd/cleaner/cleaner.go
```

## Configurations
Configurations are loaded from the environment variables, if not provided, then the default values will be used.
* `LISTEN_ADDR`: The address that the server listens on, defaults to `:8000`
* `BASE_URL`: The base URL of the service, defaults to `http://localhost:8000`
* `MAX_TOKEN`: The maximum amount of tokens to be generated and stored offline, defaults to `5000`
* `MAX_URL_LEN`: The length limit of url string, defaults to `1024`
* `MAX_ALIVE_DURATION`: The maximum duration the record stays in database, defaults to `8760h`
* `MONGODB_URI`: The uri to connect to mongodb, defaults to `mongodb://localhost:27017`
* `DATABASE`: The database to use, defaults to `shorturl`
* `RECORD_COLLECTION`: The collection to store records, defaults to `records`
* `MEMCACHED_ADDRS`: The addresses of memcached services, separated by comma, defaults to `localhost:11211`

## API Endpoints
### GET `/{id}`
### POST `/api/v1/urls`
Example Request Body (application/json)
```json
{
    "url": "https://www.dcard.tw/f",
    "expireAt": "2021-02-08T09:20:41Z"
}
```
Response
```json
{
    "id": "8RjJue",
    "shortUrl": "http://localhost:8000/8RjJue"
}
```

## Database, Cache and Packages Used

### github.com/labstack/echo
Go package `echo` is adopted as the web framework for its simplicity, extensibility and high performance. Other golang web framework packages like `gin` were also considered, but I have experiences with developing with `echo`, so it's adopted in final.

### MongoDB
Considering the use case, we basically only need to do key-value storing and searching instead of complex relation quering. So in this case, MongoDB seems to be a good choice. Also, MongoDB is well-known for its distributed architecture and high scalability, while SQL-based database like PostgreSQL is more diffcult to scale horizontally. Therefore MongoDB is adopted as the database.

### go.mongodb.org/mongo-driver/mongo
Go package `mongo` is the MongoDB Driver API package maintained by the MongoDB official, we use it to connect and operate on our database.

### Memcached
Memcached is a cache service that can be distributed deployed, in this project Memcached is adopted as the cache for speeding up the resolve time for frequently requested shorturls.

### github.com/bradfitz/gomemcache
Go package `gomemcache` is used as the interface to interact with Memcached server. Since there is no official package released, the most imported package is choosed.
