# Stage 1: Rust 検索エンジンのビルド
FROM rust:1.77 AS builder-rust
WORKDIR /app
COPY rust-search/ ./rust-search/
RUN cd rust-search && cargo build --release

# Stage 2: Go サーバのビルド
FROM golang:1.22-alpine AS builder-go
WORKDIR /app
COPY go-server/ ./go-server/
# builder-rust から Rust バイナリを go-server にコピー
COPY --from=builder-rust /app/rust-search/target/release/rust-search ./go-server/rust-search
RUN cd go-server && go build -o rustysearch ./cmd

# Stage 3: 実行用イメージ
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder-go /app/go-server/rustysearch ./
COPY --from=builder-go /app/go-server/rust-search ./
EXPOSE 8080
ENTRYPOINT ["./rustysearch"]
