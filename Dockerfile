# Use the official Go image as the build environment
FROM golang:1.21 AS build-env

# Install git and necessary certificates
RUN apt-get update && apt-get install -y git ca-certificates

# Define environment variables for GitLab
ARG GITLAB_USER
ARG GITLAB_TOKEN

# Set GOPRIVATE to allow access to private modules
RUN go env -w GOPRIVATE=gitlab.com/chicong291002/*

# Create a working directory
WORKDIR /app

# Copy go.mod and go.sum into the working directory
COPY go.* ./

# Create a .netrc file to store credentials
RUN echo "machine gitlab.com login $GITLAB_USER password $GITLAB_TOKEN" > ~/.netrc

# Download Go modules
RUN go mod download

# Copy the application source code into the image
COPY . .

# Clean up and update Go modules
RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags musl -o /api main.go

# Production stage
FROM golang:1.21

# Install necessary packages for the application
RUN apt-get update && apt-get install -y ca-certificates tzdata curl pkg-config && \
    rm -rf /var/cache/apt/*

# Copy the built binary from the previous stage
COPY --from=build-env /api /

# Expose the application on port 8080 (or the port you want)
EXPOSE 8080

# Command to run the application when the container starts
CMD ["/api"]
