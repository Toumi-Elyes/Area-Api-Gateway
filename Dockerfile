FROM golang:1.19-alpine

WORKDIR /api-gateway

COPY ./ent ent
COPY go.mod .

RUN go mod tidy

RUN go generate ./ent

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD go run .
