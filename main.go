package main

import (
	"fmt"
	"log"
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
)

// Messages for async browser operations
type pageLoadedMsg struct {
	title   string
	content string
	url     string
}

type errorMsg struct {
	err error
}

type engineInitializedMsg struct{
	engine *BrowserEngine
}

type engineInitializationFailedMsg struct {
	err error
}

// Main model for our TUI browser
type model struct {
	// Browser engine
	engine *BrowserEngine

	// UI state
	mode        viewMode
	urlInput    textinput.Model
	currentURL  string
	pageTitle   string
	pageContent string
	history     []string
	historyIdx  int

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

	return model{
		mode:               urlInputMode,
		urlInput:           ti,
		currentURL:         "",
		history:            make([]string, 0),
		historyIdx:         -1,
		width:              80,
		height:             24,
		engineReady:        false,
		engineInitializing: true,
	}
}

// Initialize the browser
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.initBrowser,
	)
}

// Initialize headless browser
func (m *model) initBrowser() tea.Msg {
	engine, err := NewBrowserEngine()
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

		title, content, finalURL, err := m.engine.Navigate(url)
		if err != nil {
			return errorMsg{err: err}
		}

		return pageLoadedMsg{
			title:   title,
			content: content,
			url:     finalURL,
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
				// Scroll down (simplified - would need better implementation)
				// In a full implementation, we'd track scroll position
			case "k", "up":
				// Scroll up
			}
		}

	case pageLoadedMsg:
		m.loading = false
		m.pageTitle = msg.title
		m.pageContent = msg.content
		m.currentURL = msg.url
		m.err = nil

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
		// Show page content with word wrapping
		wrapped := wordWrap(m.pageContent, m.width-6)
		// Limit to screen height
		lines := strings.Split(wrapped, "\n")
		maxLines := m.height - 10 // Reserve space for header and help
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			lines = append(lines, helpStyle.Render("[More content below - full scrolling not implemented in POC]"))
		}
		s.WriteString(contentStyle.Render(strings.Join(lines, "\n")))
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

	// Help text at bottom
	s.WriteString("\n\n")
	if m.mode == browsingMode {
		help := []string{
			"[l/Ctrl+L] URL bar",
			"[b] Back",
			"[f] Forward",
			"[r] Reload",
			"[q/Ctrl+C] Quit",
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