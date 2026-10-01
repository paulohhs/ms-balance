FROM golang:1.25

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/balancecore ./cmd/balancecore

EXPOSE 3003

CMD ["balancecore"]
