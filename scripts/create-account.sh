#!/bin/bash

BASE_URL="${BASE_URL:-http://localhost:8080}"

read -rp "Nome: " name
read -rp "Email: " email

response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/accounts" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"${name}\", \"email\": \"${email}\"}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo ""
  echo "Conta criada com sucesso!"
  echo "$body" | python3 -m json.tool
else
  echo ""
  echo "Erro ao criar conta (HTTP ${http_code})"
  echo "$body"
fi
