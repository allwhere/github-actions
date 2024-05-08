FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /app
COPY go.* main.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o issue-to-fixversion .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/issue-to-fixversion .
RUN ls -la
RUN chmod +x /app/issue-to-fixversion
CMD ["/app/issue-to-fixversion"]