#!/bin/bash

# Production deployment script for blog-backend

# Set secure JWT secret (generate a strong random secret)
export JWT_SECRET="$(openssl rand -base64 32)"

# Production environment variables
export GIN_MODE=release
export DATABASE_URL="./blog.db"
export PORT=8081

# SMTP Email Configuration - PLEASE CONFIGURE THESE!
# Gmail SMTP settings (recommended)
export SMTP_HOST=${SMTP_HOST:-"smtp.gmail.com"}
export SMTP_PORT=${SMTP_PORT:-"587"}
export SMTP_USER=${SMTP_USER:-""}
export SMTP_PASS=${SMTP_PASS:-""}
export FROM_EMAIL=${FROM_EMAIL:-""}

# Create .env file for permanent storage
cat > .env << EOF
JWT_SECRET=$JWT_SECRET
GIN_MODE=release
DATABASE_URL=./blog.db
PORT=8081
SMTP_HOST=$SMTP_HOST
SMTP_PORT=$SMTP_PORT
SMTP_USER=$SMTP_USER
SMTP_PASS=$SMTP_PASS
FROM_EMAIL=$FROM_EMAIL
EOF

echo "🔐 JWT_SECRET generated and saved to .env file"
echo "📝 JWT_SECRET: $JWT_SECRET"

# Check SMTP configuration
if [ -z "$SMTP_USER" ] || [ -z "$SMTP_PASS" ]; then
    echo "⚠️  WARNING: SMTP credentials not configured!"
    echo "📧 Email notifications to newsletter subscribers will not work."
    echo "🔧 To enable email notifications, set these environment variables:"
    echo "   export SMTP_USER=\"your-email@gmail.com\""
    echo "   export SMTP_PASS=\"your-app-password\""
    echo "   export FROM_EMAIL=\"your-email@gmail.com\""
    echo ""
    echo "📋 For Gmail, you need to use App Password:"
    echo "   https://myaccount.google.com/apppasswords"
    echo ""
else
    echo "✅ SMTP configuration found - email notifications enabled"
fi

# Build the application
echo "🔨 Building application..."
go build -o blog-backend .

# Make sure the binary is executable
chmod +x blog-backend

# Start the application
echo "🚀 Starting application with production settings..."
./blog-backend 