#!/bin/bash

# エラーが発生したらスクリプトを停止
set -e

echo "RustySearch MVPの起動スクリプト"
echo "==============================="

# 現在のディレクトリを基準にパスを設定
PROJECT_ROOT=$(pwd)
RUST_BIN="$PROJECT_ROOT/rust-search/target/release/fixer"
GO_SERVER_DIR="$PROJECT_ROOT/go-server/cmd"
PORT=8080

# プロセスの状態を確認
echo "0. プロセスの状態を確認中..."
PROCESS_EXISTS=false

# 8080ポートを使用しているプロセスを確認
if command -v lsof >/dev/null 2>&1; then
  EXISTING_PID=$(lsof -ti:${PORT} 2>/dev/null)
  if [ ! -z "$EXISTING_PID" ]; then
    PROCESS_EXISTS=true
    echo "ポート ${PORT} で実行中のプロセス ($EXISTING_PID) を終了します"
    kill -9 $EXISTING_PID || true
  fi
fi

# go runプロセスの確認
GO_PID=$(pgrep -f "go run main.go" 2>/dev/null)
if [ ! -z "$GO_PID" ]; then
  PROCESS_EXISTS=true
  echo "go run プロセス ($GO_PID) を終了します"
  pkill -f "go run main.go" || true
fi

# プロセスが存在した場合のみスリープ
if [ "$PROCESS_EXISTS" = true ]; then
  echo "プロセス終了後、1秒待機します..."
  sleep 1
else
  echo "実行中のプロセスはありません。続行します..."
fi

echo "1. Rustバイナリのビルド..."
cd "$PROJECT_ROOT/rust-search"
cargo build --release

echo "2. Rustバイナリの確認: $RUST_BIN"
if [ ! -f "$RUST_BIN" ]; then
  echo "エラー: Rustバイナリが見つかりません"
  exit 1
fi

echo "3. Goサーバーの起動..."
cd "$GO_SERVER_DIR"
SEARCH_ROOT="$RUST_BIN" go run main.go
