package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Browser states
type viewMode int

const (
	urlInputMode viewMode = iota
	browsingMode
	linkSelectionMode
)

// Messages for async browser operations
type pageLoadedMsg struct {
	title   string
	content string
	url     string
	links   []Link
}

type errorMsg struct {
	err error
}

type engineInitializedMsg struct {
	engine Engine
}

type engineInitializationFailedMsg struct {
	err error
}

// Main model for our TUI browser
type model struct {
	// Browser engine
	engine Engine

	// UI state
	mode        viewMode
	urlInput    textinput.Model
	currentURL  string
	pageTitle   string
	pageContent string
	history     []string
	historyIdx  int

	// Links
	links            []Link
	linkSearchInput  textinput.Model
	selectedLinkIdx  int
	linkScrollOffset int

	// Scrolling
	contentScrollOffset int
	contentLines        []string

	// Display
	width  int
	height int
	err    error

	// Loading state
	loading bool

	// Engine state
	engineReady        bool
	engineInitializing bool
}

// Initialize the model
func initialModel() model {
	// Create text input for URL
	ti := textinput.New()
	ti.Placeholder = "Enter URL (e.g., https://example.com)"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	// Create text input for link search
	linkSearch := textinput.New()
	linkSearch.Placeholder = "Search links or enter number..."
	linkSearch.CharLimit = 100
	linkSearch.Width = 80

	return model{
		mode:                urlInputMode,
		urlInput:            ti,
		linkSearchInput:     linkSearch,
		currentURL:          "",
		history:             make([]string, 0),
		historyIdx:          -1,
		links:               make([]Link, 0),
		selectedLinkIdx:     0,
		linkScrollOffset:    0,
		contentScrollOffset: 0,
		contentLines:        make([]string, 0),
		width:               80,
		height:              24,
		engineReady:         false,
		engineInitializing:  true,
	}
}

// Initialize the browser
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.initBrowser,
	)
}

// Initialize browser engine
func (m *model) initBrowser() tea.Msg {
	// Determine which engine to use
	// Default to goquery (works everywhere), but allow rod via env var
	engineType := os.Getenv("ODYSSEUS_ENGINE")
	if engineType == "" {
		engineType = "goquery" // Default
	}

	var engine Engine
	var err error

	switch engineType {
	case "rod":
		engine, err = NewRodEngine()
		if err != nil {
			// Fall back to goquery if rod fails
			log.Printf("Rod engine failed to initialize, falling back to goquery: %v", err)
			engine, err = NewGoqueryEngine()
		}
	case "goquery":
		engine, err = NewGoqueryEngine()
	default:
		// Unknown engine type, default to goquery
		log.Printf("Unknown engine type '%s', using goquery", engineType)
		engine, err = NewGoqueryEngine()
	}

	if err != nil {
		return engineInitializationFailedMsg{err: err}
	}
	return engineInitializedMsg{engine: engine}
}

// Navigate to URL
func (m *model) navigate(url string) tea.Cmd {
	return func() tea.Msg {
		// Defensive check - this should never happen with proper state management
		if m.engine == nil {
			return errorMsg{err: fmt.Errorf("browser engine is nil - this is a bug")}
		}

		title, content, finalURL, links, err := m.engine.Navigate(url)
		if err != nil {
			return errorMsg{err: err}
		}

		return pageLoadedMsg{
			title:   title,
			content: content,
			url:     finalURL,
			links:   links,
		}
	}
}


// Update handles messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.urlInput.Width = min(msg.Width-4, 100)

		// Re-wrap content if we have any
		if m.pageContent != "" {
			wrapped := wordWrap(m.pageContent, m.width-6)
			m.contentLines = strings.Split(wrapped, "\n")

			// Clamp scroll offset if content got shorter
			maxScroll := len(m.contentLines) - (m.height - 10)
			if maxScroll < 0 {
				maxScroll = 0
			}
			if m.contentScrollOffset > maxScroll {
				m.contentScrollOffset = maxScroll
			}
		}

		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case urlInputMode:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEnter:
				url := m.urlInput.Value()
				if url != "" {
					if !m.engineReady {
						if m.engineInitializing {
							m.err = fmt.Errorf("browser engine is still initializing, please wait")
						} else {
							m.err = fmt.Errorf("browser engine failed to initialize")
						}
						return m, nil
					}
					m.loading = true
					m.mode = browsingMode
					return m, m.navigate(url)
				}
			case tea.KeyEsc:
				if m.currentURL != "" {
					m.mode = browsingMode
				}
				return m, nil
			default:
				m.urlInput, cmd = m.urlInput.Update(msg)
				return m, cmd
			}

		case browsingMode:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+l", "l":
				// Open URL bar
				m.mode = urlInputMode
				m.urlInput.SetValue(m.currentURL)
				m.urlInput.Focus()
				return m, textinput.Blink
			case "/", "g":
				// Open link selection mode
				if len(m.links) > 0 {
					m.mode = linkSelectionMode
					m.selectedLinkIdx = 0
					m.linkScrollOffset = 0
					m.linkSearchInput.SetValue("")
					m.linkSearchInput.Focus()
					return m, textinput.Blink
				}
			case "b":
				// Go back in history
				if m.historyIdx > 0 && m.engineReady {
					m.historyIdx--
					url := m.history[m.historyIdx]
					m.loading = true
					return m, m.navigate(url)
				}
			case "f":
				// Go forward in history
				if m.historyIdx < len(m.history)-1 && m.engineReady {
					m.historyIdx++
					url := m.history[m.historyIdx]
					m.loading = true
					return m, m.navigate(url)
				}
			case "r":
				// Reload current page
				if m.currentURL != "" && m.engineReady {
					m.loading = true
					return m, m.navigate(m.currentURL)
				}
			case "j", "down":
				// Scroll down by 1 line
				m.scrollDown(1)
			case "k", "up":
				// Scroll up by 1 line
				m.scrollUp(1)
			case "ctrl+d":
				// Scroll down by half page
				halfPage := (m.height - 10) / 2
				if halfPage < 1 {
					halfPage = 1
				}
				m.scrollDown(halfPage)
			case "ctrl+u":
				// Scroll up by half page
				halfPage := (m.height - 10) / 2
				if halfPage < 1 {
					halfPage = 1
				}
				m.scrollUp(halfPage)
			case "ctrl+f", "pgdown":
				// Scroll down by full page
				fullPage := m.height - 10
				if fullPage < 1 {
					fullPage = 1
				}
				m.scrollDown(fullPage)
			case "ctrl+b", "pgup":
				// Scroll up by full page
				fullPage := m.height - 10
				if fullPage < 1 {
					fullPage = 1
				}
				m.scrollUp(fullPage)
			case "home":
				// Jump to top
				m.contentScrollOffset = 0
			case "end":
				// Jump to bottom
				maxScroll := len(m.contentLines) - (m.height - 10)
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.contentScrollOffset = maxScroll
			}

		case linkSelectionMode:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEsc:
				// Return to browsing mode
				m.mode = browsingMode
				m.linkSearchInput.Blur()
				return m, nil
			case tea.KeyEnter:
				// Navigate to selected link
				if m.selectedLinkIdx >= 0 && m.selectedLinkIdx < len(m.links) {
					link := m.links[m.selectedLinkIdx]
					m.mode = browsingMode
					m.loading = true
					m.linkSearchInput.Blur()
					return m, m.navigate(link.URL)
				}
			case tea.KeyUp, tea.KeyCtrlK:
				if m.selectedLinkIdx > 0 {
					m.selectedLinkIdx--
					// Adjust scroll offset if needed
					if m.selectedLinkIdx < m.linkScrollOffset {
						m.linkScrollOffset = m.selectedLinkIdx
					}
				}
				return m, nil
			case tea.KeyDown, tea.KeyCtrlJ:
				if m.selectedLinkIdx < len(m.links)-1 {
					m.selectedLinkIdx++
					// Adjust scroll offset if needed
					maxVisible := m.height - 10
					if m.selectedLinkIdx >= m.linkScrollOffset+maxVisible {
						m.linkScrollOffset = m.selectedLinkIdx - maxVisible + 1
					}
				}
				return m, nil
			default:
				// Handle number keys (0-9) for quick link selection
				if len(msg.String()) == 1 {
					char := msg.String()[0]
					if char >= '0' && char <= '9' {
						// Build number from input
						currentVal := m.linkSearchInput.Value()
						m.linkSearchInput.SetValue(currentVal + string(char))
						m.linkSearchInput, cmd = m.linkSearchInput.Update(msg)
						return m, cmd
					}
				}
				// Update search input
				m.linkSearchInput, cmd = m.linkSearchInput.Update(msg)
				// Try to parse as number for quick selection
				if val := m.linkSearchInput.Value(); val != "" {
					var num int
					if _, err := fmt.Sscanf(val, "%d", &num); err == nil && num > 0 && num <= len(m.links) {
						m.selectedLinkIdx = num - 1
					}
				}
				return m, cmd
			}
		}

	case pageLoadedMsg:
		m.loading = false
		m.pageTitle = msg.title
		m.pageContent = msg.content
		m.currentURL = msg.url
		m.links = msg.links
		m.selectedLinkIdx = 0
		m.linkScrollOffset = 0
		m.contentScrollOffset = 0
		m.err = nil

		// Prepare content lines for scrolling
		wrapped := wordWrap(msg.content, m.width-6)
		m.contentLines = strings.Split(wrapped, "\n")

		// Update history
		if m.historyIdx == -1 || m.history[m.historyIdx] != msg.url {
			// Trim forward history if we navigated from middle
			if m.historyIdx < len(m.history)-1 {
				m.history = m.history[:m.historyIdx+1]
			}
			m.history = append(m.history, msg.url)
			m.historyIdx = len(m.history) - 1
		}

		return m, nil

	case errorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case engineInitializedMsg:
		m.engine = msg.engine
		m.engineReady = true
		m.engineInitializing = false
		m.err = nil
		return m, nil

	case engineInitializationFailedMsg:
		m.engineReady = false
		m.engineInitializing = false
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

// scrollDown scrolls content down by n lines (with bounds checking)
func (m *model) scrollDown(lines int) {
	if len(m.contentLines) == 0 {
		return
	}

	maxScroll := len(m.contentLines) - (m.height - 10)
	if maxScroll < 0 {
		maxScroll = 0
	}

	m.contentScrollOffset += lines
	if m.contentScrollOffset > maxScroll {
		m.contentScrollOffset = maxScroll
	}
}

// scrollUp scrolls content up by n lines (with bounds checking)
func (m *model) scrollUp(lines int) {
	m.contentScrollOffset -= lines
	if m.contentScrollOffset < 0 {
		m.contentScrollOffset = 0
	}
}

// View renders the UI
func (m model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	contentStyle := lipgloss.NewStyle().
		Width(m.width - 4).
		Padding(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	var s strings.Builder

	// Header
	s.WriteString(titleStyle.Render("🌍 Odysseus Terminal Browser"))
	s.WriteString("\n")

	// URL bar or current URL display
	if m.mode == urlInputMode {
		s.WriteString("\n")

		// Show different input states based on engine status
		if m.engineInitializing {
			disabledInput := m.urlInput
			disabledStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Faint(true)
			s.WriteString(disabledStyle.Render(disabledInput.View()))
			s.WriteString("\n\n")
			s.WriteString(helpStyle.Render("⏳ Initializing browser engine, please wait..."))
		} else if !m.engineReady {
			disabledInput := m.urlInput
			disabledStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Faint(true)
			s.WriteString(disabledStyle.Render(disabledInput.View()))
			s.WriteString("\n\n")
			s.WriteString(errorStyle.Render("❌ Browser engine failed to initialize"))
		} else {
			s.WriteString(m.urlInput.View())
			s.WriteString("\n\n")
			s.WriteString(helpStyle.Render("Press Enter to navigate, Esc to cancel"))
		}
	} else {
		if m.currentURL != "" {
			s.WriteString(urlStyle.Render(fmt.Sprintf("📍 %s", m.currentURL)))
			s.WriteString("\n")
			if m.pageTitle != "" {
				s.WriteString(fmt.Sprintf("📄 %s", m.pageTitle))
			}
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")

	// Content area
	if m.engineInitializing && m.mode != urlInputMode {
		s.WriteString(contentStyle.Render("⏳ Initializing browser engine..."))
	} else if m.loading {
		s.WriteString(contentStyle.Render("⏳ Loading..."))
	} else if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
	} else if m.pageContent != "" {
		// Show page content with scrolling
		maxLines := m.height - 10 // Reserve space for header and help
		totalLines := len(m.contentLines)

		// Calculate visible slice
		startLine := m.contentScrollOffset
		endLine := m.contentScrollOffset + maxLines

		if endLine > totalLines {
			endLine = totalLines
		}

		// Get visible lines
		var visibleLines []string
		if startLine < totalLines {
			visibleLines = m.contentLines[startLine:endLine]
		}

		s.WriteString(contentStyle.Render(strings.Join(visibleLines, "\n")))

		// Show scroll indicators
		if totalLines > maxLines {
			s.WriteString("\n\n")

			// Calculate scroll percentage
			scrollPercent := 0
			if totalLines > maxLines {
				scrollPercent = (m.contentScrollOffset * 100) / (totalLines - maxLines)
			}

			// Show position indicator
			endLineNum := endLine
			if endLineNum > totalLines {
				endLineNum = totalLines
			}

			scrollInfo := fmt.Sprintf("Lines %d-%d of %d (%d%%)",
				startLine+1, endLineNum, totalLines, scrollPercent)

			// Add visual indicators
			if m.contentScrollOffset > 0 {
				scrollInfo = "▲ " + scrollInfo
			}
			if endLine < totalLines {
				scrollInfo = scrollInfo + " ▼"
			}

			s.WriteString(helpStyle.Render(scrollInfo))
		}
	} else if m.mode == browsingMode {
		if !m.engineReady {
			if m.engineInitializing {
				s.WriteString(helpStyle.Render("Browser engine initializing..."))
			} else {
				s.WriteString(helpStyle.Render("Browser engine failed to initialize."))
			}
		} else {
			s.WriteString(helpStyle.Render("No content loaded. Press 'l' to enter a URL."))
		}
	}

	// Link selection mode
	if m.mode == linkSelectionMode {
		s.WriteString("\n")
		s.WriteString(titleStyle.Render(fmt.Sprintf("📎 Links (%d found)", len(m.links))))
		s.WriteString("\n\n")

		// Show search input
		s.WriteString(m.linkSearchInput.View())
		s.WriteString("\n\n")

		// Calculate how many links we can show
		maxVisible := m.height - 15
		if maxVisible < 5 {
			maxVisible = 5
		}

		// Show links with selection
		endIdx := m.linkScrollOffset + maxVisible
		if endIdx > len(m.links) {
			endIdx = len(m.links)
		}

		for i := m.linkScrollOffset; i < endIdx; i++ {
			link := m.links[i]
			linkNum := i + 1

			// Style for selected vs unselected
			var linkLine string
			if i == m.selectedLinkIdx {
				selectedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("205")).
					Background(lipgloss.Color("235")).
					Bold(true)
				linkLine = selectedStyle.Render(fmt.Sprintf("→ [%d] %s", linkNum, link.Text))
			} else {
				linkLine = fmt.Sprintf("  [%d] %s", linkNum, link.Text)
			}

			s.WriteString(linkLine)
			s.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(m.links) > maxVisible {
			scrollInfo := fmt.Sprintf("  (Showing %d-%d of %d links)",
				m.linkScrollOffset+1, endIdx, len(m.links))
			s.WriteString("\n")
			s.WriteString(helpStyle.Render(scrollInfo))
		}
	}

	// Help text at bottom
	s.WriteString("\n\n")
	if m.mode == browsingMode {
		help := []string{
			"[l] URL",
			"[/,g] Links",
			"[j/k] Scroll",
			"[b/f] Back/Fwd",
			"[r] Reload",
			"[q] Quit",
		}
		s.WriteString(helpStyle.Render(strings.Join(help, " • ")))
	} else if m.mode == linkSelectionMode {
		help := []string{
			"[↑↓] Navigate",
			"[0-9] Quick select",
			"[Enter] Open",
			"[Esc] Cancel",
		}
		s.WriteString(helpStyle.Render(strings.Join(help, " • ")))
	}

	return s.String()
}

// Simple word wrap function
func wordWrap(text string, width int) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) <= width {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		// Wrap long lines
		words := strings.Fields(line)
		currentLine := ""

		for _, word := range words {
			if len(currentLine)+len(word)+1 > width {
				if currentLine != "" {
					result.WriteString(currentLine)
					result.WriteString("\n")
				}
				currentLine = word
			} else {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			}
		}

		if currentLine != "" {
			result.WriteString(currentLine)
			result.WriteString("\n")
		}
	}

	return strings.TrimSuffix(result.String(), "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Create and run the TUI application
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup browser resources
	if model, ok := finalModel.(model); ok && model.engine != nil {
		model.engine.Close()
	}
}