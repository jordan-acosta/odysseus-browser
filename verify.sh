#!/bin/bash

echo "Verifying Odysseus Terminal Browser Build..."
echo "============================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi
echo "✅ Go is installed: $(go version)"

# Check if binary was built
if [ -f "./odysseus" ]; then
    echo "✅ Odysseus binary exists"
    echo "   Size: $(ls -lh odysseus | awk '{print $5}')"
else
    echo "❌ Odysseus binary not found"
    exit 1
fi

# Check dependencies
echo ""
echo "Dependencies installed:"
go list -m all | grep -E "(charmbracelet|rod)" | while read -r line; do
    echo "  ✓ $line"
done

echo ""
echo "Build verification complete!"
echo ""
echo "To run the browser:"
echo "  ./odysseus"
echo ""
echo "Or use the demo script:"
echo "  ./demo.sh"