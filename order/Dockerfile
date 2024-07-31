#FROM golang:1.19 AS builder
#WORKDIR /usr/src/app
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go
#
#FROM scratch
#COPY --from=builder /usr/src/app/order ./order
#CMD ["./order"]

FROM golang:1.22

# Set the current working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=99
ENV GRPC_GO_LOG_SEVERITY_LEVEL=info
ENV GRPC_XDS_BOOTSTRAP=/var/run/config/xds_bootstrap.json
RUN go mod tidy
# Build the Go app
RUN go build -o order  ./cmd/main.go

# Command to run the executable
CMD ["./order"]