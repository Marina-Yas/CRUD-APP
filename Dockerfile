# syntax=docker/dockerfile:1
FROM golang:1.17 As builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crud-app .

FROM alpine
WORKDIR /app
ENV PORT 8080
COPY --from=builder /app/crud-app ./
CMD ["./crud-app"]