FROM golang:1.18-rc
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build cmd/server/main.go

CMD ./main