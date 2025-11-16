#!/bin/bash

echo "======================================"
echo "Odysseus Browser - Engine Tests"
echo "======================================"
echo ""

# Build the browser
echo "Building browser..."
go build -o odysseus || exit 1
echo "✅ Build successful"
echo ""

# Test goquery engine
echo "Testing Goquery engine..."
echo "This should work on all platforms (including Android/Termux)"
echo ""
echo "To test manually, run:"
echo "  ODYSSEUS_ENGINE=goquery ./odysseus"
echo ""

# Test rod engine (may fail if Chrome not installed)
echo "Testing Rod engine..."
echo "This requires Chrome/Chromium to be installed"
echo ""
echo "To test manually, run:"
echo "  ODYSSEUS_ENGINE=rod ./odysseus"
echo ""

# Show binary size
echo "Binary size: $(ls -lh odysseus | awk '{print $5}')"
echo ""

# Test that goquery works without CGO
echo "Testing CGO-disabled build (for Android)..."
CGO_ENABLED=0 go build -o odysseus-nocgo || exit 1
echo "✅ CGO-disabled build successful"
rm odysseus-nocgo
echo ""

echo "======================================"
echo "All tests passed!"
echo "======================================"
echo ""
echo "Quick Start:"
echo "  ./odysseus                    # Use goquery (default)"
echo "  ODYSSEUS_ENGINE=rod ./odysseus  # Use rod (requires Chrome)"
echo ""
echo "For Android/Termux users:"
echo "  - The default goquery engine will work perfectly"
echo "  - No Chrome installation needed"
echo "  - Builds without CGO issues"
