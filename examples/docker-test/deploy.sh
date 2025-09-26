echo "🚀 Starting deployment with MyTerraform..."

# Step 1: Build
echo "📦 Building MyTerraform..."
go build -o derraform ./cmd/app

# Step 2: Validate configuration
echo "🔍 Validating configuration..."
./myterraform plan examples/web-app.tf

# Step 3: Apply
echo "⚡ Applying configuration..."
./myterraform apply examples/web-app.tf

# Step 4: Verify
echo "✅ Verifying deployment..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo "🎉 Deployment completed!"