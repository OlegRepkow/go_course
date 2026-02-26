#!/bin/bash

echo "=== Testing Chat Service ==="
echo ""

BASE_URL="http://localhost:8080"

echo "1. Sign Up"
SIGNUP=$(curl -s -X POST $BASE_URL/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}')
echo "Response: $SIGNUP"
echo ""

echo "2. Sign In"
SIGNIN=$(curl -s -X POST $BASE_URL/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}')
echo "Response: $SIGNIN"
TOKEN=$(echo $SIGNIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"
echo ""

if [ -z "$TOKEN" ]; then
  echo "Failed to get token"
  exit 1
fi

echo "3. Send Message"
SEND=$(curl -s -X POST $BASE_URL/channel/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"text":"Hello from test!"}')
echo "Response: $SEND"
echo ""

echo "4. Get History"
HISTORY=$(curl -s -X GET $BASE_URL/channel/history \
  -H "Authorization: Bearer $TOKEN")
echo "Response: $HISTORY"
echo ""

echo "=== Test Complete ==="
