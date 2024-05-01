# Use a light version of Go
FROM golang:1.16-alpine as builder

WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o jira-release-creator .

# Use a minimal base image to run the binary
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/jira-release-creator .

# Command to run the executable
CMD ["./jira-release-creator"]