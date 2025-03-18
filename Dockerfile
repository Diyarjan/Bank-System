FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 GOOS=linux  go build -o ./bank
RUN go build -o bank .

EXPOSE 8080
RUN chmod a+x bank
ENTRYPOINT ["./bank"]

#FROM builder AS  run-test-stage
#RUN go test -v ./...

#FROM gcr.io/distroless/base-debian11 AS build-release-stage

#WORKDIR /app
#COPY --from=builder /build/bank ./bank
#RUN chmod +x ./bank  # Даем права на выполнение

#ENV PORT=8080

#EXPOSE 8080
#USER nonroot:nonroot
#ENTRYPOINT ["/app/bank"]

#docker build -t new-bank:1.0 .
#docker run -p 8080:8080 new-bank:1.0
