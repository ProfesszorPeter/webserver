FROM golang

WORKDIR /app

COPY go.mod ./

EXPOSE 8000

RUN go mod download

COPY . .

RUN go build -o /app/main /app/main.go

VOLUME webserver-data:/var/lib/webserver

ENTRYPOINT ["/app/main"]
