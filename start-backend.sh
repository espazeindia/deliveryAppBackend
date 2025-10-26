#!/bin/bash

# Espaze Delivery App Backend Startup Script

echo "🚀 Starting Espaze Delivery App Backend..."
echo ""

# Check if MongoDB is running
if ! pgrep -x mongod > /dev/null; then
    echo "⚠️  MongoDB is not running. Starting MongoDB..."
    brew services start mongodb-community
    sleep 3
fi

# Kill any existing backend processes
echo "🔄 Stopping any existing backend processes..."
pkill -9 -f "delivery-backend" 2>/dev/null || true
lsof -ti:8081 | xargs kill -9 2>/dev/null || true
sleep 2

# Set environment variables
export PORT=8081
export MONGO_URI="mongodb://localhost:27017/espaze_delivery"
export JWT_SECRET="espaze-delivery-secret-key-2025-production"

# Get local IP
LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)

# Build and start backend
echo "📦 Building backend..."
go build -o /tmp/delivery-backend main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful"
    echo ""
    echo "🌐 Starting backend server..."
    nohup /tmp/delivery-backend > /tmp/delivery-backend.log 2>&1 &
    BACKEND_PID=$!
    echo "Backend PID: $BACKEND_PID"
    sleep 4
    
    # Test health endpoint
    HEALTH=$(curl -s http://localhost:8081/health)
    if echo "$HEALTH" | grep -q "ok"; then
        echo "✅ Backend is running successfully!"
        echo ""
        echo "📡 Backend URLs:"
        echo "   Local:   http://localhost:8081"
        echo "   Network: http://$LOCAL_IP:8081"
        echo ""
        echo "📋 API Endpoints:"
        echo "   Health:  http://localhost:8081/health"
        echo "   Login:   POST http://localhost:8081/api/v1/delivery/login"
        echo "   Orders:  GET  http://localhost:8081/api/v1/delivery/orders/active"
        echo ""
        echo "📱 Connect your iPhone to: http://$LOCAL_IP:8081"
        echo ""
        echo "📝 View logs: tail -f /tmp/delivery-backend.log"
        echo "🛑 Stop backend: pkill delivery-backend"
    else
        echo "❌ Backend failed to start. Check logs:"
        tail -20 /tmp/delivery-backend.log
        exit 1
    fi
else
    echo "❌ Build failed"
    exit 1
fi

