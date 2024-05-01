FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /app
COPY go.* main.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o jira-release-creator .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/jira-release-creator .
RUN ls -la
RUN chmod +x jira-release-creator
CMD ["./jira-release-creator"]