#!/bin/sh

API_URL=${API_URL:-"http://127.0.0.1:3000"}
UUID=${1:-"$(uuidgen)"}
NAME=${2:-"John Doe"}
EMAIL=${3:-"john@example.com"}
DOB=${4:-"1990-01-01T00:00:00Z"}

echo "Creating user with UUID: ${UUID}"

curl -v -X "POST" "${API_URL}/save" \
  -H "Content-Type: application/json" \
  -d "{
    \"external_id\": \"${UUID}\",
    \"name\": \"${NAME}\",
    \"email\": \"${EMAIL}\",
    \"date_of_birth\": \"${DOB}\"
  }"
