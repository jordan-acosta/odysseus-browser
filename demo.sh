#!/bin/bash

echo "======================================"
echo "Odysseus Terminal Browser - Demo"
echo "======================================"
echo ""
echo "This proof-of-concept demonstrates a terminal web browser"
echo "that uses Chromium for rendering and Charm for the TUI."
echo ""
echo "Key Features:"
echo "  - Headless Chromium browser integration via go-rod"
echo "  - Terminal UI built with Charm's Bubble Tea framework"
echo "  - URL navigation with text input"
echo "  - Basic browsing history (back/forward)"
echo "  - Text content extraction and display"
echo ""
echo "Controls:"
echo "  - Enter URL mode: Press 'l' or Ctrl+L"
echo "  - Navigate: Type URL and press Enter"
echo "  - Go Back: Press 'b'"
echo "  - Go Forward: Press 'f'"
echo "  - Reload: Press 'r'"
echo "  - Quit: Press 'q' or Ctrl+C"
echo ""
echo "Starting browser in 3 seconds..."
sleep 3

./odysseus