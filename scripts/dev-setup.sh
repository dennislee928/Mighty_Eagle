#!/bin/bash

# Aegis Trust Ecosystem - Development Setup Script

echo "🚀 Setting up Aegis Trust Ecosystem for development..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start database services
echo "📦 Starting PostgreSQL and Redis..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are healthy
if docker-compose ps | grep -q "healthy"; then
    echo "✅ Database services are ready!"
else
    echo "⚠️  Services may still be starting up. Check with: docker-compose ps"
fi

# Install dependencies if not already installed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Build shared packages
echo "🔧 Building shared packages..."
npm run build --workspace=@aegis/shared

echo "🎉 Setup complete! You can now run:"
echo "   npm run dev    # Start all development servers"
echo "   docker-compose ps    # Check service status"
echo ""
echo "Services will be available at:"
echo "   Web App: http://localhost:3000"
echo "   API: http://localhost:8080"
echo "   PostgreSQL: localhost:5432"
echo "   Redis: localhost:6379"