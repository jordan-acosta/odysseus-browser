#!/bin/bash

# Session Start Hook for Watch Collection Expo App
# This script automatically starts the Expo dev server in tunnel mode
# when a Claude Code session begins.

echo "🚀 Starting Expo dev server in tunnel mode..."
echo "📱 Open Expo Go on your phone and scan the QR code below"
echo ""

# Navigate to the watch-collection directory
cd "$(dirname "$0")/../.." || exit 1

# Start Expo in tunnel mode (runs in background)
# The tunnel mode allows your phone to connect from anywhere
npm run start:tunnel &

# Store the process ID
EXPO_PID=$!
echo "Expo dev server started with PID: $EXPO_PID"
echo ""
echo "💡 Tips:"
echo "  - Scan the QR code with your Expo Go app"
echo "  - Changes you make will hot-reload on your phone"
echo "  - Use 'npm run start:tunnel' to restart if needed"
echo ""
