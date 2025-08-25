FROM golang:1.25-alpine as builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o listener .
CMD ["/app/listener"]