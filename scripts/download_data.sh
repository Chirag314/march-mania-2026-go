#!/usr/bin/env bash
set -euo pipefail

COMP="march-machine-learning-mania-2026"
OUT="data"

mkdir -p "$OUT"

# Prefer new token-based auth if present
if [[ -n "${KAGGLE_API_TOKEN:-}" ]]; then
  echo "Using KAGGLE_API_TOKEN auth"
elif [[ -n "${KAGGLE_USERNAME:-}" && -n "${KAGGLE_KEY:-}" ]]; then
  echo "Using KAGGLE_USERNAME/KAGGLE_KEY auth"
else
  echo "Missing credentials. Set either:"
  echo "  - KAGGLE_API_TOKEN"
  echo "  OR"
  echo "  - KAGGLE_USERNAME and KAGGLE_KEY"
  exit 1
fi

kaggle competitions download -c "$COMP" -p "$OUT" --force
unzip -o "$OUT"/*.zip -d "$OUT"
