FROM golang:1.20 AS builder
WORKDIR /src
RUN go install github.com/cosmtrek/air@latest
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app /app
EXPOSE 8080

ENTRYPOINT ["/app"]
