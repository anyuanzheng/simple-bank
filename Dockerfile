# Build stage
FROM golang:1.19 AS builder
WORKDIR /app
COPY . .
# CGO has to be disabled for alpine
RUN export CGO_ENABLED=0 && go build -o main main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]
