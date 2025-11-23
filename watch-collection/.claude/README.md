# Watch Collection App - Claude Code Setup

This Expo app is configured for development using Claude Code for the web with live preview on your phone via Expo Go.

## Quick Start

### First Time Setup

1. **Install Expo Go on your phone**
   - iOS: Download from App Store
   - Android: Download from Google Play Store

2. **Start the dev server**
   - The session start hook automatically launches the server when Claude Code starts
   - Or manually run: `cd watch-collection && npm run start:tunnel`

3. **Connect from your phone**
   - Scan the QR code displayed in the terminal with Expo Go
   - The app will load on your phone
   - Any changes made in Claude Code will hot-reload automatically

## Development Workflow

### Making Changes

1. Edit files in the `watch-collection` directory
2. Save your changes
3. Watch them appear instantly on your phone (hot reload)

### Project Structure

```
watch-collection/
├── app/              # Expo Router pages
├── components/       # React components
├── constants/        # App constants (colors, etc.)
├── hooks/            # Custom React hooks
├── assets/           # Images, fonts, etc.
└── .claude/          # Claude Code configuration
```

### Key Commands

```bash
# Start dev server with tunnel (for phone access)
npm run start:tunnel

# Start regular dev server
npm start

# Run linter
npm run lint

# Reset project (if needed)
npm run reset-project
```

## How Tunnel Mode Works

- Tunnel mode creates a publicly accessible URL for your dev server
- This allows your phone to connect from any network
- Perfect for Claude Code on the web (cloud-based development)
- The URL is temporary and changes each session

## Troubleshooting

### QR Code not scanning?
- Make sure both devices are connected to internet
- Try typing the URL manually in Expo Go
- Restart the dev server: `npm run start:tunnel`

### App not updating?
- Check if the dev server is still running
- Try shaking your phone and selecting "Reload"
- Check for errors in the terminal

### Session start hook not working?
- Manually start the server: `cd watch-collection && npm run start:tunnel`
- Check that the hook file is executable: `chmod +x .claude/hooks/session-start.sh`

## Migrating to New Repository

When ready to move this project to `jordan-acosta/claude-code-mobile-test`:

1. Copy the entire `watch-collection` directory contents
2. Push to the new repo
3. The `.claude` directory will configure Claude Code automatically
4. The session start hook will work in the new repo

## About This Project

This is a watch collection management app built with:
- **Expo** - React Native framework
- **Expo Router** - File-based navigation
- **TypeScript** - Type-safe code
- **Expo Go** - Live preview on phone

The app will help you track:
- Watches you own
- Important information about each watch
- Photos and details
- Collection management
