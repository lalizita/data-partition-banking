#!/bin/bash

BASE_URL="${BASE_URL:-http://localhost:8081}"

read -rp "Client ID (UUID): " client_id
read -rp "Valor: " amount

response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/transactions/credit" \
  -H "Content-Type: application/json" \
  -d "{\"client_id\": \"${client_id}\", \"amount\": ${amount}}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo ""
  echo "Transação criada com sucesso!"
  [ -n "$body" ] && echo "$body" | python3 -m json.tool 2>/dev/null || echo "$body"
else
  echo ""
  echo "Erro ao criar transação (HTTP ${http_code})"
  echo "$body"
fi
