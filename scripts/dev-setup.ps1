# Aegis Trust Ecosystem Development Setup Script for Windows

Write-Host "Setting up Aegis Trust Ecosystem development environment..." -ForegroundColor Green

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    Write-Host "Docker is available" -ForegroundColor Green
} catch {
    Write-Host "Docker is not available. Please install Docker Desktop and ensure it's running." -ForegroundColor Red
    exit 1
}

# Check if Node.js is installed
Write-Host "Checking Node.js..." -ForegroundColor Yellow
try {
    node --version | Out-Null
    Write-Host "Node.js is available" -ForegroundColor Green
} catch {
    Write-Host "Node.js is not available. Please install Node.js 18 or later." -ForegroundColor Red
    exit 1
}

# Check if Go is installed
Write-Host "Checking Go..." -ForegroundColor Yellow
try {
    go version | Out-Null
    Write-Host "Go is available" -ForegroundColor Green
} catch {
    Write-Host "Go is not available. Please install Go 1.21 or later." -ForegroundColor Red
    exit 1
}

# Install Node.js dependencies
Write-Host "Installing Node.js dependencies..." -ForegroundColor Yellow
npm install

# Start database services
Write-Host "Starting PostgreSQL and Redis..." -ForegroundColor Yellow
docker-compose up -d postgres redis

# Wait for services to be ready
Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check if services are healthy
Write-Host "Checking service health..." -ForegroundColor Yellow
$postgresHealth = docker-compose ps postgres --format json | ConvertFrom-Json | Select-Object -ExpandProperty Health
$redisHealth = docker-compose ps redis --format json | ConvertFrom-Json | Select-Object -ExpandProperty Health

if ($postgresHealth -eq "healthy" -and $redisHealth -eq "healthy") {
    Write-Host "All services are healthy!" -ForegroundColor Green
} else {
    Write-Host "Services may not be fully ready. Check docker-compose logs for details." -ForegroundColor Yellow
}

# Copy environment files
Write-Host "Setting up environment files..." -ForegroundColor Yellow
if (!(Test-Path "apps/web/.env.local")) {
    Copy-Item "apps/web/.env.local.example" "apps/web/.env.local"
    Write-Host "Created apps/web/.env.local from example" -ForegroundColor Green
}

Write-Host "Development environment setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Run 'npm run dev' to start all development servers" -ForegroundColor White
Write-Host "2. Visit http://localhost:3000 for the web app" -ForegroundColor White
Write-Host "3. Visit http://localhost:8080/health for the API health check" -ForegroundColor White
Write-Host ""
Write-Host "Useful commands:" -ForegroundColor Cyan
Write-Host "- npm run db:up    # Start database services" -ForegroundColor White
Write-Host "- npm run db:down  # Stop database services" -ForegroundColor White
Write-Host "- npm run db:reset # Reset database with fresh data" -ForegroundColor White