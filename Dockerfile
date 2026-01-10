FROM golang:1.24.11-bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    gcc \
    g++ \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8084

RUN CGO_ENABLED=1 go build -o delivery-service .

CMD ["./delivery-service"]
