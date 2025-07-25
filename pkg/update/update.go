package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/nizar0x1f/termup/pkg/version"
)

const (
	githubAPIURL = "https://api.github.com/repos/nizar0x1f/termup/releases/latest"

	checkInterval = 24 * time.Hour
)

type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	ReleaseURL     string
	ReleaseNotes   string
}

func CheckForUpdates() (*UpdateInfo, error) {
	currentVersion := version.Get().Version

	if currentVersion == "dev" || currentVersion == "unknown" {
		return &UpdateInfo{
			Available:      false,
			CurrentVersion: currentVersion,
		}, nil
	}

	latest, err := getLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	if latest.Prerelease || latest.Draft {
		return &UpdateInfo{
			Available:      false,
			CurrentVersion: currentVersion,
		}, nil
	}

	latestVersion := strings.TrimPrefix(latest.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	isNewer, err := isVersionNewer(latestVersion, currentVersionClean)
	if err != nil {
		return nil, fmt.Errorf("failed to compare versions: %w", err)
	}

	return &UpdateInfo{
		Available:      isNewer,
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		ReleaseURL:     latest.HTMLURL,
		ReleaseNotes:   latest.Body,
	}, nil
}

func getLatestRelease() (*Release, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(githubAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func isVersionNewer(newVersion, currentVersion string) (bool, error) {

	newParts, err := parseVersion(newVersion)
	if err != nil {
		return false, err
	}

	currentParts, err := parseVersion(currentVersion)
	if err != nil {
		return false, err
	}

	for i := 0; i < 3; i++ {
		if newParts[i] > currentParts[i] {
			return true, nil
		}
		if newParts[i] < currentParts[i] {
			return false, nil
		}
	}

	return false, nil
}

func parseVersion(version string) ([3]int, error) {
	var parts [3]int

	version = strings.TrimPrefix(version, "v")

	versionParts := strings.Split(version, ".")
	if len(versionParts) != 3 {
		return parts, fmt.Errorf("invalid version format: %s", version)
	}

	for i, part := range versionParts {
		var num int
		if _, err := fmt.Sscanf(part, "%d", &num); err != nil {
			return parts, fmt.Errorf("invalid version part: %s", part)
		}
		parts[i] = num
	}

	return parts, nil
}

func GetUpdateCommand() string {
	return "go install github.com/nizar0x1f/termup/cmd/upl@latest"
}

func FormatUpdateMessage(info *UpdateInfo) string {
	if !info.Available {
		return ""
	}

	msg := "A new version of TermUp (S3 compatible filesharing) is available!\n"
	msg += "Current version: " + info.CurrentVersion + "\n"
	msg += "Latest version:  " + info.LatestVersion + "\n"
	msg += "\nTo update, run:\n"
	msg += "  " + GetUpdateCommand() + "\n"
	msg += "\nRelease notes: " + info.ReleaseURL + "\n"

	return msg
}

func ShouldCheckForUpdates(lastCheck time.Time) bool {
	return time.Since(lastCheck) > checkInterval
}

func PerformSelfUpdate() error {

	if !isGoAvailable() {
		return fmt.Errorf("Go is not installed or not in PATH. Please install Go or update manually using: %s", GetUpdateCommand())
	}

	cmd := exec.Command("go", "install", "github.com/nizar0x1f/termup/cmd/upl@latest")

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go install failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func isGoAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func GetBinaryPath() (string, error) {
	return os.Executable()
}

func CanSelfUpdate() (bool, string) {
	if !isGoAvailable() {
		return false, "Go is not installed or not in PATH"
	}

	execPath, err := GetBinaryPath()
	if err != nil {
		return false, "Cannot determine executable path: " + err.Error()
	}

	if runtime.GOOS == "windows" {
		return true, "Self-update available (binary: " + execPath + ")"
	}

	return true, "Self-update available (binary: " + execPath + ")"
}
