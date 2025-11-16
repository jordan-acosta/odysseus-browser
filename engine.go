package main

// Link represents a clickable link on the page
type Link struct {
	Text string // Display text
	URL  string // Absolute URL
}

// Engine defines the interface for browser engines
// This allows switching between goquery (pure Go) and rod (headless Chrome)
type Engine interface {
	// Navigate loads a URL and returns page info
	Navigate(url string) (title, content, finalURL string, links []Link, err error)

	// Reload refreshes the current page
	Reload() error

	// Close cleans up engine resources
	Close()
}

// EngineType represents the type of browser engine
type EngineType string

const (
	EngineTypeGoquery EngineType = "goquery" // Pure Go HTML parser (default, works everywhere)
	EngineTypeRod     EngineType = "rod"     // Headless Chrome (supports JS, requires Chrome)
)
