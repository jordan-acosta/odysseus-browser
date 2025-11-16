# Odysseus World-Wide Wanderer

A proof-of-concept terminal web browser built in Golang with support for multiple rendering engines.

## Features

- 🌍 **Dual Engine Support**:
  - **Goquery** (default): Pure Go HTML parser - works everywhere including Android/Termux
  - **Rod**: Headless Chrome with JavaScript support
- 🔗 **Interactive Navigation**: Select and navigate links with keyboard
- 📱 **Cross-Platform**: Works on Linux, macOS, Windows, Android (Termux), and ARM devices
- ⚡ **Fast & Lightweight**: ~17MB binary, minimal dependencies
- 🎨 **Beautiful TUI**: Built with Charm's Bubble Tea framework

## Quick Start

```bash
# Build (works on any platform)
go build -o odysseus

# Run with default engine (goquery - works everywhere)
./odysseus

# Or use headless Chrome (requires Chrome/Chromium)
ODYSSEUS_ENGINE=rod ./odysseus
```

## Building for Android/Termux

```bash
# In Termux
pkg install golang
git clone <repo>
cd odysseus-browser
go build -o odysseus
./odysseus
```

No CGO issues, no Chrome required!

## Controls

- `l` or `Ctrl+L` - Open URL bar
- `/` or `g` - Browse links on current page
- `↑`/`↓` - Navigate links
- `0-9` - Quick select link by number
- `Enter` - Navigate to selected link
- `b`/`f` - Back/Forward in history
- `r` - Reload page
- `q` or `Ctrl+C` - Quit

## Engine Documentation

See [ENGINES.md](ENGINES.md) for detailed information about the two browser engines.

## Work in Progress

This is a proof-of-concept demonstrating terminal-based web browsing with pluggable rendering engines.
