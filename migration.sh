#!/usr/bin/env bash
set -euo pipefail

DB_URL="postgres://abdullah:@localhost:5432/gator"

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 {up|down}"
  exit 1
fi

case "$1" in
  up)
    goose -dir sql/schema postgres "$DB_URL" up
    ;;
  down)
    goose -dir sql/schema postgres "$DB_URL" down
    ;;
  *)
    echo "Unknown command: $1"
    echo "Usage: $0 {up|down}"
    exit 1
    ;;
esac