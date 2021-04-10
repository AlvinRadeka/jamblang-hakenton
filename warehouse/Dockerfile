FROM golang:1.16-alpine AS builder

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o ./build/app ./cmd

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/build/app /app/build/app
COPY --from=builder /app/configs/config.yaml /app/configs/config.yaml

CMD "/app/build/app"