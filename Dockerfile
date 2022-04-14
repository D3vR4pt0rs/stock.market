FROM golang:latest AS builder
ENV PROJECT_PATH=/app/market
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . ${PROJECT_PATH}
WORKDIR ${PROJECT_PATH}
RUN go build cmd/market/main.go

FROM golang:alpine
WORKDIR /app/cmd/market
COPY --from=builder /app/market/main .
EXPOSE 1337
CMD ["./main"]