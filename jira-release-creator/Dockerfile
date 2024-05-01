# Build stage
FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o jira-release-creator .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/jira-release-creator .

# Command to run the executable
ENTRYPOINT ["./jira-release-creator"]
