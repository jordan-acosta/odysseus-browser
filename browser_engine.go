package main

import (
	"context"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// BrowserEngine handles all browser operations
type BrowserEngine struct {
	browser *rod.Browser
	page    *rod.Page
}

// NewBrowserEngine creates a new browser engine instance
func NewBrowserEngine() (*BrowserEngine, error) {
	// Try to find Chrome/Chromium on the system
	path, exists := launcher.LookPath()
	var l string

	if exists {
		// Use system Chrome/Chromium if available
		l = launcher.New().
			Bin(path).
			Headless(true).
			NoSandbox(true). // Required for some environments
			MustLaunch()
	} else {
		// Download and use embedded browser if needed
		l = launcher.New().
			Headless(true).
			NoSandbox(true).
			MustLaunch()
	}

	browser := rod.New().
		ControlURL(l).
		MustConnect()

	page := browser.MustPage()

	// Set a reasonable viewport size for text extraction
	page.MustSetViewport(1024, 768, 1, false)

	// Set user agent to avoid bot detection
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	})

	return &BrowserEngine{
		browser: browser,
		page:    page,
	}, nil
}

// Navigate loads a URL and returns page info
func (b *BrowserEngine) Navigate(url string) (title, content, finalURL string, err error) {
	// Ensure URL has protocol
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Navigate with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = b.page.Context(ctx).Navigate(url)
	if err != nil {
		return "", "", "", err
	}

	// Wait for page to load
	err = b.page.WaitLoad()
	if err != nil {
		return "", "", "", err
	}

	// Get page title
	titleElem, err := b.page.Element("title")
	if err == nil && titleElem != nil {
		title, _ = titleElem.Text()
	}
	if title == "" {
		title = "Untitled"
	}

	// Extract text content
	content = b.extractContent()

	// Get final URL (after redirects)
	finalURL = b.page.MustInfo().URL

	return title, content, finalURL, nil
}

// extractContent gets readable text from the page
func (b *BrowserEngine) extractContent() string {
	// Use JavaScript to extract and clean text content
	bodyText := b.page.MustEval(`() => {
		// Remove script and style elements
		const scripts = document.querySelectorAll('script, style, noscript, iframe');
		scripts.forEach(el => el.remove());

		// Function to extract text from an element
		function getText(element) {
			// Skip hidden elements
			const style = window.getComputedStyle(element);
			if (style.display === 'none' || style.visibility === 'hidden') {
				return '';
			}

			// Try to find main content areas first
			const main = element.querySelector('main, article, [role="main"], #content, .content, .main-content');
			if (main) {
				return main.innerText || main.textContent || '';
			}

			// Try to find and exclude navigation, footer, etc.
			const contentElements = element.querySelectorAll('p, h1, h2, h3, h4, h5, h6, li, td, th, pre, blockquote');
			let text = [];

			contentElements.forEach(el => {
				// Skip if element is inside nav, footer, aside, or header
				if (!el.closest('nav, footer, aside, header')) {
					const content = (el.innerText || el.textContent || '').trim();
					if (content && content.length > 0) {
						text.push(content);
					}
				}
			});

			if (text.length > 0) {
				return text.join('\n\n');
			}

			// Fallback to body text
			return element.innerText || element.textContent || '';
		}

		const body = document.body;
		if (!body) return "No content available";

		return getText(body);
	}`).String()

	// Clean up the text
	lines := strings.Split(bodyText, "\n")
	var cleanLines []string
	prevEmpty := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip multiple empty lines
		if line == "" {
			if !prevEmpty {
				cleanLines = append(cleanLines, "")
				prevEmpty = true
			}
			continue
		}

		prevEmpty = false
		cleanLines = append(cleanLines, line)
	}

	content := strings.Join(cleanLines, "\n")

	// Limit content length for terminal display
	const maxLen = 10000
	if len(content) > maxLen {
		// Try to cut at a word boundary
		cutPoint := maxLen
		for i := maxLen; i > maxLen-100 && i >= 0; i-- {
			if content[i] == ' ' || content[i] == '\n' {
				cutPoint = i
				break
			}
		}
		content = content[:cutPoint] + "\n\n[Content truncated for display...]"
	}

	return content
}

// Reload refreshes the current page
func (b *BrowserEngine) Reload() error {
	return b.page.Reload()
}

// Close cleans up browser resources
func (b *BrowserEngine) Close() {
	if b.page != nil {
		b.page.Close()
	}
	if b.browser != nil {
		b.browser.Close()
	}
}