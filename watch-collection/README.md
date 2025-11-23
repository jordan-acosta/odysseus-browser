# Watch Collection Manager 🕰️

A mobile app for managing and tracking your watch collection. Built with Expo and React Native, designed for easy iteration using Claude Code on the web with live preview on your phone.

## Features (Planned)

- 📸 Add watches with photos and details
- 🏷️ Track important information (brand, model, purchase date, value)
- 📊 View your collection at a glance
- 🔍 Search and filter watches
- 💰 Track collection value
- 📝 Add notes and maintenance records

## Getting Started

### Prerequisites

- **Expo Go app** on your phone ([iOS](https://apps.apple.com/app/expo-go/id982107779) | [Android](https://play.google.com/store/apps/details?id=host.exp.exponent))
- Node.js installed (if running locally)

### Development with Claude Code

This project is optimized for development using Claude Code for the web:

1. **Automatic Setup**: The session start hook automatically launches the dev server
2. **Scan QR Code**: Open Expo Go and scan the QR code shown in the terminal
3. **Live Updates**: Changes sync instantly to your phone

See `.claude/README.md` for detailed Claude Code setup instructions.

### Manual Development

If you want to run this locally:

```bash
# Install dependencies
npm install

# Start dev server (tunnel mode - for phone access)
npm run start:tunnel

# Or use the helper script
./dev-server.sh
```

## Project Structure

```
watch-collection/
├── app/              # Expo Router pages (file-based routing)
│   ├── (tabs)/       # Tab navigation screens
│   ├── +not-found.tsx
│   └── _layout.tsx
├── components/       # Reusable React components
│   ├── ui/           # UI components (buttons, cards, etc.)
│   └── ParallaxScrollView.tsx
├── constants/        # App constants
│   └── Colors.ts     # Color scheme
├── hooks/            # Custom React hooks
├── assets/           # Images, fonts, and static files
├── .claude/          # Claude Code configuration
│   ├── hooks/        # Session hooks
│   └── README.md     # Claude Code documentation
└── scripts/          # Utility scripts
```

## Tech Stack

- **Expo SDK 54** - React Native framework with managed workflow
- **Expo Router** - File-based navigation system
- **TypeScript** - Type-safe JavaScript
- **React Native** - Cross-platform mobile development
- **Expo Go** - Development client for live preview

## Development Commands

```bash
# Start development server
npm start               # Local mode
npm run start:tunnel    # Tunnel mode (recommended for remote access)

# Platform-specific
npm run android         # Open on Android emulator/device
npm run ios            # Open on iOS simulator (macOS only)
npm run web            # Open in web browser

# Code quality
npm run lint           # Run ESLint

# Utilities
npm run reset-project  # Reset to fresh state
./dev-server.sh        # Interactive dev server helper
```

## Migrating to Production Repository

This project is currently in the `odysseus-browser/watch-collection` subdirectory.

To move to `jordan-acosta/claude-code-mobile-test`:

```bash
# Copy everything from watch-collection directory
cp -r watch-collection/* /path/to/new/repo/

# The .claude configuration will work automatically in the new repo
```

## Contributing

This is a personal project for learning mobile development with Claude Code. Feel free to fork and adapt for your own use!

## License

Private - Personal Use Only

## Next Steps

1. Design the data model for watch information
2. Implement watch list screen
3. Add watch detail view
4. Implement camera integration for photos
5. Add local storage (AsyncStorage or SQLite)
6. Build search and filter functionality
7. Add data export/backup features

---

Built with ❤️ using Claude Code and Expo
