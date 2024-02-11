# ビルドステージ
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subscribe ./cmd/subscribe/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o genkeys ./cmd/genkeys/main.go

# 実行ステージ
FROM alpine:latest

# タイムゾーンを指定するために入れる
RUN apk --update-cache add tzdata

COPY --from=builder /app/subscribe .
COPY --from=builder /app/genkeys .

CMD ["./subscribe"]
