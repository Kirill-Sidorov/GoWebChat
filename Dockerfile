FROM golang:1.22.5-alpine3.19 AS builder
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o .

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=builder /resources /resources
COPY --from=builder .env webchat /
EXPOSE 8080
ENTRYPOINT ["/webchat"]