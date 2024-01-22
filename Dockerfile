FROM golang:1.21 as BUILD

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-list-service

CMD ["/app/todo-list-service"]