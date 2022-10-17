FROM golang:1.17-alpine as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o server cmd/daily/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server .

EXPOSE 4000

CMD [ "./server" ]