# Use official Golang image
FROM golang:1.22.2

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o api .

# Expose the port
EXPOSE 3000

# Run the executable
CMD ["./api"]
