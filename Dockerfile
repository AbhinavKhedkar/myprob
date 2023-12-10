# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set execute permission for the main binary
RUN chmod +x main

# Expose port 8000
EXPOSE 8000

# Command to run the application
CMD ["./main"]
