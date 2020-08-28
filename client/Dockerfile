FROM golang:1.12

WORKDIR /app

COPY go.* ./
RUN go get

COPY *.go ./
RUN go build -o mqtt-client

CMD ["./mqtt-client"]