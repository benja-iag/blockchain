FROM golang:alpine AS build
WORKDIR /go/src/blockchain
COPY . .
RUN go build -o main main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/blockchain/ .
CMD ["./main", "createblockchain", "a"]
