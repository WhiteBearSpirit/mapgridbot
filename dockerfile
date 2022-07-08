FROM golang:1.17
WORKDIR /app
COPY *.go ./
COPY go.* ./
RUN go build -o /gridbot
CMD ["/gridbot"]
