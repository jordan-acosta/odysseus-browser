# Migration Guide: Moving Watch Collection to New Repository

This guide explains how to move the `watch-collection` Expo app from this repository to the `jordan-acosta/claude-code-mobile-test` repository.

## Quick Migration Steps

### Option 1: Manual Copy (Recommended)

1. **Clone the new repository locally**:
   ```bash
   git clone https://github.com/jordan-acosta/claude-code-mobile-test.git
   cd claude-code-mobile-test
   ```

2. **Copy the watch-collection directory contents**:
   ```bash
   # From the odysseus-browser repository
   cp -r /path/to/odysseus-browser/watch-collection/* .

   # Or copy individual items:
   cp -r /path/to/odysseus-browser/watch-collection/.* .
   cp -r /path/to/odysseus-browser/watch-collection/* .
   ```

3. **Commit and push to the new repo**:
   ```bash
   git add .
   git commit -m "Initial commit: Watch Collection Expo app with Claude Code setup"
   git push origin main
   ```

### Option 2: Using Git (Preserve History)

If you want to preserve git history:

```bash
# From odysseus-browser directory
cd watch-collection

# Initialize git if needed
git init

# Add remote for new repo
git remote add new-origin https://github.com/jordan-acosta/claude-code-mobile-test.git

# Commit and push
git add .
git commit -m "Initial commit: Watch Collection Expo app"
git push new-origin main
```

## What Gets Migrated

The following will be copied to the new repository:

### Application Code
- `app/` - Expo Router pages and navigation
- `components/` - React components
- `constants/` - App constants (colors, etc.)
- `hooks/` - Custom React hooks
- `assets/` - Images, fonts, and static files
- `scripts/` - Utility scripts

### Configuration Files
- `package.json` - Dependencies and scripts (includes tunnel mode)
- `app.json` - Expo configuration
- `tsconfig.json` - TypeScript configuration
- `eslint.config.js` - ESLint configuration
- `.gitignore` - Git ignore rules

### Claude Code Setup
- `.claude/hooks/session-start.sh` - Auto-starts dev server
- `.claude/README.md` - Claude Code documentation

### Development Tools
- `dev-server.sh` - Interactive dev server helper
- `README.md` - Project documentation

### Dependencies
- `node_modules/` - Will be regenerated with `npm install`
- `package-lock.json` - Lock file for dependencies

## Post-Migration Checklist

After copying files to the new repository:

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Verify the setup**:
   ```bash
   # Start dev server in tunnel mode
   npm run start:tunnel
   ```

3. **Test Claude Code integration**:
   - Open the new repo in Claude Code for the web
   - The session start hook should automatically start the dev server
   - Scan the QR code with Expo Go on your phone
   - Verify hot reload works

4. **Update any hardcoded paths**:
   - Check if any files reference the old repository path
   - Update them to reflect the new location

## Troubleshooting

### Session Hook Not Running

If the session start hook doesn't run automatically:

```bash
# Make sure the hook is executable
chmod +x .claude/hooks/session-start.sh

# Manually run the hook to test
./.claude/hooks/session-start.sh
```

### Dependencies Missing

If you see module errors:

```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

### QR Code Not Displaying

If the QR code doesn't show:

```bash
# Try installing the Expo CLI globally
npm install -g @expo/cli

# Or use npx
npx expo start --tunnel
```

## Verification

Once migrated, verify everything works:

- [ ] `npm install` runs successfully
- [ ] `npm run start:tunnel` starts the dev server
- [ ] QR code is displayed
- [ ] Expo Go on phone can scan and connect
- [ ] Hot reload works when editing files
- [ ] `.claude/` configuration is present
- [ ] Session start hook works in Claude Code

## Next Steps After Migration

1. **Update the repository description** on GitHub
2. **Add collaborators** if needed
3. **Configure branch protection** rules
4. **Start developing** your watch collection app!

## Need Help?

- Check `.claude/README.md` for Claude Code specific help
- Check `watch-collection/README.md` for app-specific documentation
- Review Expo documentation: https://docs.expo.dev

---

After successful migration, you can safely delete the `watch-collection` directory from the `odysseus-browser` repository.
