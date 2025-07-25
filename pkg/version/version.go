package version

import (
	"fmt"
	"runtime"
	"time"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: GoVersion,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func (i Info) String() string {
	if i.Version == "dev" {
		return fmt.Sprintf("termup %s - S3 compatible filesharing (commit: %s, built: %s, go: %s, platform: %s)",
			i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
	}
	return fmt.Sprintf("termup v%s - S3 compatible filesharing (commit: %s, built: %s, go: %s, platform: %s)",
		i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
}

func (i Info) Short() string {
	if i.Version == "dev" {
		return "termup " + i.Version
	}
	return "termup v" + i.Version
}

func (i Info) IsRelease() bool {
	return i.Version != "dev" && i.Version != "unknown"
}

func (i Info) ParseDate() (time.Time, error) {
	if i.Date == "unknown" {
		return time.Time{}, fmt.Errorf("build date unknown")
	}
	return time.Parse(time.RFC3339, i.Date)
}
