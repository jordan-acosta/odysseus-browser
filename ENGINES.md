# Browser Engine Backends

Odysseus supports two browser engine backends:

## 1. Goquery (Default) ⭐

**Pure Go HTML parser** - works everywhere, no dependencies

### Pros
- ✅ Works on ALL platforms (Linux, macOS, Windows, Android, ARM, etc.)
- ✅ No CGO required - compiles anywhere Go runs
- ✅ No external dependencies (Chrome/Chromium not needed)
- ✅ Fast and lightweight
- ✅ Perfect for Android/Termux
- ✅ Works great for static sites and many modern sites

### Cons
- ❌ No JavaScript execution
- ❌ Can't handle Single Page Applications (SPAs) that render client-side
- ❌ Dynamic content won't be visible

### Best For
- News sites (HN, Reddit, news outlets)
- Blogs and documentation
- Wikipedia
- GitHub (mostly works)
- Stack Overflow
- Any site that works without JavaScript

### Usage
```bash
# Default - just run the browser
./odysseus

# Or explicitly set goquery
ODYSSEUS_ENGINE=goquery ./odysseus
```

## 2. Rod (Headless Chrome)

**Full browser engine** - supports JavaScript execution

### Pros
- ✅ Executes JavaScript
- ✅ Handles SPAs and dynamic content
- ✅ Works with modern web apps
- ✅ Renders pages exactly like a real browser

### Cons
- ❌ Requires Chrome/Chromium installed
- ❌ Uses CGO on some platforms (build issues)
- ❌ Doesn't work well on Android/Termux
- ❌ Heavier resource usage
- ❌ Slower to start

### Best For
- JavaScript-heavy sites
- Single Page Applications
- Modern web apps
- Sites that don't work with goquery

### Usage
```bash
ODYSSEUS_ENGINE=rod ./odysseus
```

## Building for Different Platforms

### Android/Termux (Use Goquery)
```bash
# Install build tools if needed
pkg install golang

# Build with goquery (default)
go build -o odysseus

# Or build without CGO explicitly
CGO_ENABLED=0 go build -o odysseus
```

### Linux/macOS (Either Engine)
```bash
# For goquery only (smaller binary, works everywhere)
CGO_ENABLED=0 go build -tags=goquery -o odysseus

# For both engines (default)
go build -o odysseus
```

### Cross-compilation
```bash
# For Raspberry Pi (ARM)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o odysseus-arm64

# For Android (ARM64)
GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -o odysseus-android
```

## Automatic Fallback

If you set `ODYSSEUS_ENGINE=rod` but Rod fails to initialize (e.g., Chrome not found), Odysseus will automatically fall back to goquery mode. This ensures the browser always works.

## Engine Comparison

| Feature | Goquery | Rod |
|---------|---------|-----|
| JavaScript Support | ❌ | ✅ |
| Works on Android | ✅ | ❌ |
| Build Complexity | Simple | Requires Chrome |
| Binary Size | ~16MB | ~16MB |
| Memory Usage | Low | Higher |
| Startup Time | Fast | Slower |
| Static Sites | ✅ | ✅ |
| Dynamic Sites | Partial | ✅ |

## Recommended Setup

- **Development/Desktop**: Use default (goquery) for most sites, switch to rod for JS-heavy sites
- **Android/Termux**: Use goquery only
- **Raspberry Pi**: Use goquery only
- **Server/Cloud**: Either engine works, goquery preferred for simplicity

## Testing Engines

Try these sites to see the difference:

**Work well with both:**
- news.ycombinator.com
- old.reddit.com
- wikipedia.org

**Better with goquery (simpler):**
- Most documentation sites
- Blogs
- News sites

**Require rod:**
- twitter.com/x.com
- Modern web apps
- Sites with client-side rendering
