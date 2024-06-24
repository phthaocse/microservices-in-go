#FROM golang:1.22 AS builder
#WORKDIR /usr/src/app
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go
#
#FROM scratch
#COPY --from=builder /usr/src/app/payment ./payment
#CMD ["./payment"]

# Use the official Golang image as a base image
FROM golang:1.22

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

RUN go mod tidy
# Build the Go app
RUN go build -o payment  ./cmd/main.go

# Command to run the executable
CMD ["./payment"]
