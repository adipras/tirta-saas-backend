#!/bin/bash

# Script to generate Postman collection from Huma OpenAPI

echo "🚀 Starting Tirta SaaS Backend with Huma..."

# Start the server in background
go run main.go &
SERVER_PID=$!

# Wait for server to start
echo "⏳ Waiting for server to start..."
sleep 5

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Server failed to start"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "✅ Server is running"

# Download OpenAPI JSON
echo "📥 Downloading OpenAPI specification..."
curl -s http://localhost:8080/openapi.json > openapi.json

if [ ! -s openapi.json ]; then
    echo "❌ Failed to download OpenAPI spec"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "✅ OpenAPI spec downloaded"

# Convert OpenAPI to Postman Collection using npx (requires Node.js)
echo "🔄 Converting OpenAPI to Postman collection..."

if command -v npx &> /dev/null; then
    npx openapi-to-postmanv2 -s openapi.json -o Tirta-SaaS-Backend.postman_collection.json -p
    
    if [ -f "Tirta-SaaS-Backend.postman_collection.json" ]; then
        echo "✅ Postman collection generated successfully!"
        echo "📁 File: Tirta-SaaS-Backend.postman_collection.json"
    else
        echo "⚠️  Postman collection generation failed"
        echo "   You can manually import openapi.json to Postman"
    fi
else
    echo "⚠️  Node.js not found. Please install Node.js to auto-generate Postman collection"
    echo "   Or manually import openapi.json to Postman"
    echo "   Postman: Import > Link > paste: http://localhost:8080/openapi.json"
fi

# Stop the server
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null

echo ""
echo "📚 Available files:"
echo "   - openapi.json (OpenAPI 3.1 specification)"
if [ -f "Tirta-SaaS-Backend.postman_collection.json" ]; then
    echo "   - Tirta-SaaS-Backend.postman_collection.json (Postman collection)"
fi

echo ""
echo "✨ Done! You can now:"
echo "   1. Import openapi.json to Postman, Insomnia, or any API client"
echo "   2. View interactive docs at: http://localhost:8080/docs (when server is running)"
echo "   3. Use the Postman collection for testing"
