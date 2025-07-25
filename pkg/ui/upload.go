package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UploadModel struct {
	progress   progress.Model
	spinner    spinner.Model
	filename   string
	url        string
	err        error
	done       bool
	uploading  bool
	fileSize   int64
	uploaded   int64
	startTime  time.Time
	lastUpdate time.Time
	lastBytes  int64
	speed      float64
}

var (
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	filenameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF")).
			Underline(true)

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB"))

	speedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
)

func NewUploadModel(filename string, fileSize int64) UploadModel {
	p := progress.New(progress.WithDefaultGradient())
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	now := time.Now()
	return UploadModel{
		progress:   p,
		spinner:    s,
		filename:   filename,
		fileSize:   fileSize,
		uploading:  true,
		startTime:  now,
		lastUpdate: now,
		lastBytes:  0,
		speed:      0,
	}
}

func (m UploadModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.progress.Init(),
	)
}

func (m UploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter", "q":
			if m.done {
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case UploadProgressMsg:
		now := time.Now()
		newUploaded := int64(msg)

		if !m.lastUpdate.IsZero() && newUploaded > m.lastBytes {
			timeDiff := now.Sub(m.lastUpdate).Seconds()
			if timeDiff > 0 {
				bytesDiff := newUploaded - m.lastBytes
				m.speed = float64(bytesDiff) / timeDiff
			}
		}

		m.uploaded = newUploaded
		m.lastUpdate = now
		m.lastBytes = newUploaded

		if m.fileSize > 0 {
			percent := float64(m.uploaded) / float64(m.fileSize)
			return m, m.progress.SetPercent(percent)
		}
		return m, nil

	case UploadCompleteMsg:
		m.url = string(msg)
		m.done = true
		m.uploading = false
		return m, nil

	case UploadErrorMsg:
		m.err = error(msg)
		m.done = true
		m.uploading = false
		return m, nil
	}

	return m, nil
}

func (m UploadModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("File Upload"))
	b.WriteString("\n\n")

	b.WriteString("Uploading: ")
	b.WriteString(filenameStyle.Render(m.filename))
	b.WriteString("\n\n")

	if m.uploading {

		b.WriteString(m.spinner.View())
		b.WriteString(" ")
		b.WriteString(m.progress.View())
		b.WriteString("\n\n")

		if m.fileSize > 0 {

			percent := float64(m.uploaded) / float64(m.fileSize) * 100

			elapsed := time.Since(m.startTime)

			var eta string
			if m.speed > 0 && m.uploaded > 0 {
				remaining := m.fileSize - m.uploaded
				etaSeconds := float64(remaining) / m.speed
				eta = formatDuration(time.Duration(etaSeconds * float64(time.Second)))
			} else {
				eta = "--:--"
			}

			b.WriteString(statsStyle.Render(fmt.Sprintf(
				"%s / %s (%.1f%%) %s/s ETA: %s",
				formatBytes(m.uploaded),
				formatBytes(m.fileSize),
				percent,
				speedStyle.Render(formatBytes(int64(m.speed))),
				eta,
			)))
			b.WriteString("\n")

			b.WriteString(statsStyle.Render("Elapsed: " + formatDuration(elapsed)))
		}
	} else if m.err != nil {
		b.WriteString(errorStyle.Render("✗ Upload failed"))
		b.WriteString("\n\n")
		b.WriteString("Error: " + m.err.Error())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press q or enter to exit"))
	} else {
		b.WriteString(successStyle.Render("✓ Upload successful!"))
		b.WriteString("\n\n")
		b.WriteString("URL: ")
		b.WriteString(urlStyle.Render(m.url))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press q or enter to exit"))
	}

	return b.String()
}

type UploadProgressMsg int64
type UploadCompleteMsg string
type UploadErrorMsg error

func (m *UploadModel) UpdateProgress(uploaded int64) tea.Cmd {
	return func() tea.Msg {
		return UploadProgressMsg(uploaded)
	}
}

func (m *UploadModel) CompleteUpload(url string) tea.Cmd {
	return func() tea.Msg {
		return UploadCompleteMsg(url)
	}
}

func (m *UploadModel) ErrorUpload(err error) tea.Cmd {
	return func() tea.Msg {
		return UploadErrorMsg(err)
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m UploadModel) IsDone() bool {
	return m.done
}

func (m UploadModel) GetURL() string {
	return m.url
}

func (m UploadModel) GetError() error {
	return m.err
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%02.0fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
