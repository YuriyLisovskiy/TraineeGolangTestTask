# syntax=docker/dockerfile:1
FROM golang:1.17-alpine
WORKDIR /app/
ENV CGO_ENABLED=1
RUN apk add build-base
COPY . .
RUN go mod download
RUN go build -o ./rest-api-app ./main.go

FROM alpine:3.16
WORKDIR /app/
COPY --from=0 /app/rest-api-app ./
CMD ./rest-api-app migrate && ./rest-api-app --bind 0.0.0.0:$PORT
