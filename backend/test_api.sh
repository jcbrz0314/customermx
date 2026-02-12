#!/bin/bash

echo "Testing CustomerMX API"
echo "====================="
echo ""

# Test health
echo "1. Testing /health endpoint..."
curl -s http://localhost:8080/health | jq .
echo ""

# Test login
echo "2. Testing /auth/login endpoint..."
response=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@customermx.com","password":"admin123"}')
echo "$response" | jq .
echo ""

# Extract token if login was successful
access_token=$(echo "$response" | jq -r '.data.access_token // empty')

if [ -n "$access_token" ]; then
  echo "✅ Login successful!"
  echo "Token: ${access_token:0:50}..."
  echo ""

  # Test /auth/me
  echo "3. Testing /auth/me endpoint..."
  curl -s -X GET http://localhost:8080/api/v1/auth/me \
    -H "Authorization: Bearer $access_token" | jq .
  echo ""

  # Test list brands
  echo "4. Testing /brands endpoint..."
  curl -s -X GET http://localhost:8080/api/v1/brands \
    -H "Authorization: Bearer $access_token" | jq '.data | length'
  echo " brands found"
  echo ""
else
  echo "❌ Login failed"
  echo ""
fi
