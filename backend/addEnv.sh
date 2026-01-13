#!/usr/bin/env bash
set -e

FILE="$HOME/.bashrc"
LINE='export CODELABORATE_DATABASE_URL="postgres://codelaborateuser:codelaborate@localhost:5432/codelaboratedb"'

touch "$FILE"

if ! grep -qxF "$LINE" "$FILE"; then
  echo "$LINE" >> "$FILE"
fi

