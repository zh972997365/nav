FROM golang:1.24-alpine AS builder

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum* ./
RUN go mod download

COPY backend ./backend

WORKDIR /src/backend
RUN go build -o /out/nav-backend .

FROM alpine:3.20

WORKDIR /app/backend

COPY backend/.env ./.env
COPY --from=builder /out/nav-backend ./main

RUN mkdir -p ./data

EXPOSE 8080

CMD ["./main"]
