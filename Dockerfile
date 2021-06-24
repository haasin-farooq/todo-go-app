FROM golang:1.15.0-alpine

# Maintainer information
LABEL maintainer="Haasin Farooq"

# Move to working directory /app
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the all code files into the container
COPY . .

# Build the application
RUN go build

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["./todo-go-app"]