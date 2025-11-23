#!/bin/bash

# Development Server Helper Script
# This script helps manage the Expo dev server for the Watch Collection app

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Watch Collection - Dev Server Helper   ║${NC}"
echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "$PROJECT_DIR/package.json" ]; then
    echo -e "${YELLOW}Error: package.json not found${NC}"
    echo "Please run this script from the watch-collection directory"
    exit 1
fi

# Function to show menu
show_menu() {
    echo -e "${GREEN}Available commands:${NC}"
    echo "  1) Start dev server (tunnel mode)"
    echo "  2) Start dev server (local mode)"
    echo "  3) Run linter"
    echo "  4) Install dependencies"
    echo "  5) Clear cache and restart"
    echo "  6) Show help"
    echo "  7) Exit"
    echo ""
}

# Function to start tunnel mode
start_tunnel() {
    echo -e "${BLUE}Starting Expo dev server in tunnel mode...${NC}"
    echo -e "${YELLOW}Scan the QR code with Expo Go on your phone${NC}"
    npm run start:tunnel
}

# Function to start local mode
start_local() {
    echo -e "${BLUE}Starting Expo dev server in local mode...${NC}"
    npm start
}

# Function to run linter
run_lint() {
    echo -e "${BLUE}Running linter...${NC}"
    npm run lint
}

# Function to install dependencies
install_deps() {
    echo -e "${BLUE}Installing dependencies...${NC}"
    npm install
}

# Function to clear cache
clear_cache() {
    echo -e "${BLUE}Clearing cache and restarting...${NC}"
    rm -rf .expo
    rm -rf node_modules/.cache
    echo -e "${GREEN}Cache cleared!${NC}"
    echo -e "${BLUE}Starting dev server...${NC}"
    npm start
}

# Function to show help
show_help() {
    echo -e "${GREEN}Watch Collection Dev Server Helper${NC}"
    echo ""
    echo "This script helps you manage the Expo development server."
    echo ""
    echo -e "${YELLOW}Quick Start:${NC}"
    echo "  - Choose option 1 to start in tunnel mode (recommended for phone)"
    echo "  - Scan the QR code with Expo Go app on your phone"
    echo "  - Make changes in the code and see them update live"
    echo ""
    echo -e "${YELLOW}Troubleshooting:${NC}"
    echo "  - If the app isn't updating, try option 5 (clear cache)"
    echo "  - If you see dependency errors, try option 4 (install deps)"
    echo ""
}

# If script is called with an argument, execute that command
if [ $# -eq 0 ]; then
    # Interactive mode
    while true; do
        show_menu
        read -p "Enter your choice [1-7]: " choice

        case $choice in
            1) start_tunnel ;;
            2) start_local ;;
            3) run_lint ;;
            4) install_deps ;;
            5) clear_cache ;;
            6) show_help ;;
            7) echo "Goodbye!"; exit 0 ;;
            *) echo -e "${YELLOW}Invalid option. Please try again.${NC}" ;;
        esac

        echo ""
        read -p "Press Enter to continue..."
        clear
    done
else
    # Direct command mode
    case $1 in
        tunnel) start_tunnel ;;
        local) start_local ;;
        lint) run_lint ;;
        install) install_deps ;;
        clear) clear_cache ;;
        help) show_help ;;
        *)
            echo -e "${YELLOW}Unknown command: $1${NC}"
            echo "Available commands: tunnel, local, lint, install, clear, help"
            exit 1
            ;;
    esac
fi
