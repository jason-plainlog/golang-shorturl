# build stage
FROM golang:alpine AS builder
COPY . /app
WORKDIR /app
RUN go get -u ./...
RUN go build ./cmd/server/server.go
RUN go build ./cmd/cleaner/cleaner.go

# final image
FROM golang:alpine
WORKDIR /app
COPY --from=builder /app/server ./
COPY --from=builder /app/cleaner ./

CMD ./server