#!/bin/bash

API_URL=${API_URL:-"http://127.0.0.1:3000"}
UUID=$1

if [ -z ${UUID} ]; then
  echo "Usage: $0 <external_id>"
  exit 1
fi

echo "Fetching user with UUID: $UUID"

curl -X "GET" ${API_URL}/${UUID}
