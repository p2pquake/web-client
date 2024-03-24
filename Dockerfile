FROM node:20-slim as tailwind
WORKDIR /app
COPY package.json package-lock.json /app/
RUN npm install
COPY tailwind.config.js generate.sh /app/
COPY template /app/template
RUN ./generate.sh

FROM golang:1.22-bullseye as builder
WORKDIR /go/src
COPY go.mod go.sum /go/src/
RUN go mod download
ADD . /go/src
RUN CGO_ENABLED=0 go build . && ls -l /go/src

FROM alpine:latest
WORKDIR /go
COPY --from=builder /go/src/web-client .
COPY --from=builder /go/src/static ./static
COPY --from=builder /go/src/template ./template
COPY --from=tailwind /app/static/main.css ./static/main.css
CMD ["./web-client"]
