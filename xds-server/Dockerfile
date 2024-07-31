FROM golang:alpine3.20 as build

WORKDIR /app


COPY go.mod go.sum ./
# COPY go-control-plane/go.mod go-control-plane/go.sum ./go-control-plane/


COPY . .
RUN go mod tidy
RUN go build  -o main ./cmd/
# RUN chmod +x main

FROM alpine:3.15
COPY --from=build /app/main /

ENTRYPOINT [ "./main" ]
