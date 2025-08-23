FROM golang:1.24.2-alpine as dev

WORKDIR /api

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build ./cmd/api/main.go


FROM alpine:3.19
WORKDIR /root/

COPY --from=dev /api/main .

CMD ["./main"]