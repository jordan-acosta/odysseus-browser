package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GoqueryEngine implements Engine using goquery (pure Go HTML parser)
// Works everywhere, no CGO required, perfect for Android/Termux
type GoqueryEngine struct {
	client     *http.Client
	currentURL string
	doc        *goquery.Document
}

// NewGoqueryEngine creates a new goquery-based browser engine
func NewGoqueryEngine() (*GoqueryEngine, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return &GoqueryEngine{
		client: client,
	}, nil
}

// Navigate loads a URL and extracts content
func (g *GoqueryEngine) Navigate(urlStr string) (title, content, finalURL string, links []Link, err error) {
	// Ensure URL has protocol
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	// Create request with realistic headers
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Fetch the page
	resp, err := g.client.Do(req)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Get final URL after redirects
	finalURL = resp.Request.URL.String()
	g.currentURL = finalURL

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	g.doc = doc

	// Extract title
	title = doc.Find("title").First().Text()
	if title == "" {
		title = "Untitled"
	}

	// Extract content
	content = g.extractContent(doc)

	// Extract links
	links = g.extractLinks(doc, finalURL)

	return title, content, finalURL, links, nil
}

// extractContent extracts readable text from the page
func (g *GoqueryEngine) extractContent(doc *goquery.Document) string {
	// Remove script, style, and other non-content elements
	doc.Find("script, style, noscript, iframe").Remove()

	var content strings.Builder

	// Try to find main content area first
	mainSelectors := []string{
		"main",
		"article",
		"[role='main']",
		"#content",
		".content",
		".main-content",
		"#main",
	}

	var mainContent *goquery.Selection
	for _, selector := range mainSelectors {
		mainContent = doc.Find(selector).First()
		if mainContent.Length() > 0 {
			break
		}
	}

	// If no main content found, use body
	if mainContent == nil || mainContent.Length() == 0 {
		mainContent = doc.Find("body")
	}

	// Extract text from content elements
	mainContent.Find("p, h1, h2, h3, h4, h5, h6, li, td, th, pre, blockquote, div").Each(func(i int, s *goquery.Selection) {
		// Skip if inside nav, footer, aside, or header
		if s.Closest("nav, footer, aside, header").Length() > 0 {
			return
		}

		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 0 {
			content.WriteString(text)
			content.WriteString("\n\n")
		}
	})

	result := content.String()

	// Clean up excessive whitespace
	lines := strings.Split(result, "\n")
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

	result = strings.Join(cleanLines, "\n")

	// Limit content length for terminal display
	const maxLen = 10000
	if len(result) > maxLen {
		// Try to cut at a word boundary
		cutPoint := maxLen
		for i := maxLen; i > maxLen-100 && i >= 0; i-- {
			if result[i] == ' ' || result[i] == '\n' {
				cutPoint = i
				break
			}
		}
		result = result[:cutPoint] + "\n\n[Content truncated for display...]"
	}

	return result
}

// extractLinks extracts all clickable links from the page
func (g *GoqueryEngine) extractLinks(doc *goquery.Document, baseURL string) []Link {
	var links []Link
	seen := make(map[string]bool) // Track duplicates by URL+text

	base, err := url.Parse(baseURL)
	if err != nil {
		return links
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		// Skip javascript:, mailto:, tel:, etc.
		if strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "mailto:") ||
			strings.HasPrefix(href, "tel:") ||
			strings.HasPrefix(href, "#") {
			return
		}

		// Resolve relative URLs
		linkURL, err := base.Parse(href)
		if err != nil {
			return
		}
		absoluteURL := linkURL.String()

		// Get link text
		text := strings.TrimSpace(s.Text())
		if text == "" {
			// Try to get alt text from images
			imgAlt, exists := s.Find("img").First().Attr("alt")
			if exists && imgAlt != "" {
				text = imgAlt
			} else {
				text = absoluteURL // Use URL as fallback
			}
		}

		// Truncate very long link text
		if len(text) > 100 {
			text = text[:97] + "..."
		}

		// Skip duplicates
		key := absoluteURL + "|" + text
		if seen[key] {
			return
		}
		seen[key] = true

		links = append(links, Link{
			Text: text,
			URL:  absoluteURL,
		})
	})

	return links
}

// Reload refreshes the current page
func (g *GoqueryEngine) Reload() error {
	if g.currentURL == "" {
		return fmt.Errorf("no page to reload")
	}

	_, _, _, _, err := g.Navigate(g.currentURL)
	return err
}

// Close cleans up resources (goquery doesn't need cleanup, but we implement for interface)
func (g *GoqueryEngine) Close() {
	// Nothing to clean up for goquery
	g.doc = nil
}

// fetchWithRetry attempts to fetch content with retries
func (g *GoqueryEngine) fetchWithRetry(url string, maxRetries int) (io.ReadCloser, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := g.client.Get(url)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				return resp.Body, nil
			}
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		// Wait before retry (exponential backoff)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(1<<uint(i)) * time.Second)
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
