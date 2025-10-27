# Aegis Trust Ecosystem - Development Setup Script (PowerShell)

Write-Host "🚀 Setting up Aegis Trust Ecosystem for development..." -ForegroundColor Green

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker is not running. Please start Docker and try again." -ForegroundColor Red
    exit 1
}

# Start database services
Write-Host "📦 Starting PostgreSQL and Redis..." -ForegroundColor Yellow
docker-compose up -d

# Wait for services to be ready
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check if services are healthy
$services = docker-compose ps
if ($services -match "healthy") {
    Write-Host "✅ Database services are ready!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Services may still be starting up. Check with: docker-compose ps" -ForegroundColor Yellow
}

# Install dependencies if not already installed
if (!(Test-Path "node_modules")) {
    Write-Host "📦 Installing dependencies..." -ForegroundColor Yellow
    npm install
}

# Build shared packages
Write-Host "🔧 Building shared packages..." -ForegroundColor Yellow
npm run build --workspace=@aegis/shared

Write-Host "🎉 Setup complete! You can now run:" -ForegroundColor Green
Write-Host "   npm run dev    # Start all development servers" -ForegroundColor Cyan
Write-Host "   docker-compose ps    # Check service status" -ForegroundColor Cyan
Write-Host ""
Write-Host "Services will be available at:" -ForegroundColor Green
Write-Host "   Web App: http://localhost:3000" -ForegroundColor Cyan
Write-Host "   API: http://localhost:8080" -ForegroundColor Cyan
Write-Host "   PostgreSQL: localhost:5432" -ForegroundColor Cyan
Write-Host "   Redis: localhost:6379" -ForegroundColor Cyan