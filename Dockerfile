FROM golang


COPY main.go /
COPY static/ /static

EXPOSE 8000

VOLUME webserver-data:/var/lib/webserver

RUN go get cloud.google.com/go/storage
RUN go build -o /main /main.go

ENTRYPOINT ["/main"]
