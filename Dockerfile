FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd


FROM alpine:3.14

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/myapp /myapp

ENV PORT=8080

EXPOSE 8080

CMD ["/myapp"]