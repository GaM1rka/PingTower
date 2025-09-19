FROM golang:alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/ && \
    chmod +x main

EXPOSE 8082

CMD ["./main"]