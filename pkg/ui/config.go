package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nizar0x1f/termup/pkg/config"
)

type ConfigModel struct {
	step     int
	inputs   []string
	current  string
	config   *config.Config
	finished bool
}

const (
	stepAccessKey = iota
	stepSecretKey
	stepBucket
	stepEndpoint
	stepPublicUrl
	stepComplete
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B7")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

func NewConfigModel() ConfigModel {
	return ConfigModel{
		step:   stepAccessKey,
		inputs: make([]string, 5),
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return nil
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.step == stepComplete {
				return m, m.createConfig()
			}

			m.inputs[m.step] = m.current
			m.current = ""
			m.step++
			if m.step == stepComplete {
				return m, m.createConfig()
			}

		case "backspace":
			if len(m.current) > 0 {
				m.current = m.current[:len(m.current)-1]
			}

		case "ctrl+u":

			m.current = ""

		default:

			input := msg.String()

			if len(input) == 1 {

				if input[0] >= 32 && input[0] <= 126 {
					m.current += input
				}
			} else if len(input) > 1 {

				cleaned := cleanPastedInput(input)
				if cleaned != "" {
					m.current += cleaned
				}
			}
		}

	case configCreatedMsg:
		m.config = (*config.Config)(msg)
		m.finished = true
		return m, tea.Quit
	}

	return m, nil
}

func (m ConfigModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("S3 Storage Configuration"))
	b.WriteString("\n\n")

	prompts := []string{
		"Enter Access Key ID:",
		"Enter Secret Access Key:",
		"Enter Bucket Name:",
		"Enter S3 Endpoint:",
		"Enter Public URL (press enter for default):",
	}

	if m.step < len(prompts) {
		b.WriteString(promptStyle.Render(prompts[m.step]))
		b.WriteString("\n")
		b.WriteString(inputStyle.Render(m.current + "█"))
		b.WriteString("\n\n")

		if m.step == stepPublicUrl {
			b.WriteString(helpStyle.Render("Default: https://your-bucket.s3.amazonaws.com/"))
			b.WriteString("\n\n")
		}

		b.WriteString(helpStyle.Render("Press Enter to continue, Ctrl+C to quit, Ctrl+U to clear"))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Tip: Paste works with Ctrl+V, Cmd+V, or right-click"))
	} else {
		b.WriteString("Configuration complete! Press Enter to save.")
	}

	return b.String()
}

func (m ConfigModel) createConfig() tea.Cmd {
	return func() tea.Msg {
		PublicUrl := strings.TrimSpace(m.inputs[stepPublicUrl])
		if PublicUrl == "" {
			PublicUrl = "https://your-bucket.s3.amazonaws.com/"
		}

		cfg := &config.Config{
			AccessKeyID:     strings.TrimSpace(m.inputs[stepAccessKey]),
			SecretAccessKey: strings.TrimSpace(m.inputs[stepSecretKey]),
			Bucket:          strings.TrimSpace(m.inputs[stepBucket]),
			Endpoint:        strings.TrimSpace(m.inputs[stepEndpoint]),
			PublicUrl:       PublicUrl,
		}

		return configCreatedMsg(cfg)
	}
}

type configCreatedMsg *config.Config

func (m ConfigModel) GetConfig() *config.Config {
	return m.config
}

func (m ConfigModel) IsFinished() bool {
	return m.finished
}

func cleanPastedInput(input string) string {
	cleaned := input

	cleaned = strings.TrimPrefix(cleaned, "\x1b[200~")
	cleaned = strings.TrimSuffix(cleaned, "\x1b[201~")

	if strings.HasPrefix(cleaned, "[") && strings.HasSuffix(cleaned, "]") && len(cleaned) > 2 {
		inner := cleaned[1 : len(cleaned)-1]

		shouldRemoveBrackets := false

		if strings.Contains(inner, "://") {
			shouldRemoveBrackets = true
		}

		if strings.Contains(inner, ".") {
			shouldRemoveBrackets = true
		}

		if len(inner) > 10 {
			shouldRemoveBrackets = true
		}

		if len(inner) >= 3 && (strings.Contains(strings.ToUpper(inner), "AKIA") ||
			strings.Contains(inner, "-") || strings.Contains(inner, "_")) {
			shouldRemoveBrackets = true
		}

		if len(inner) >= 3 && isValidBucketName(inner) {
			shouldRemoveBrackets = true
		}

		if shouldRemoveBrackets {
			cleaned = inner
		}
	}

	filtered := ""
	for _, r := range cleaned {
		if r >= 32 && r <= 126 {
			filtered += string(r)
		}
	}

	return filtered
}

func isValidBucketName(name string) bool {

	if len(name) < 3 || len(name) > 63 {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '.' || r == '_') {
			return false
		}
	}

	return true
}
