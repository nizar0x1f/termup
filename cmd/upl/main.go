package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nizar0x1f/termup/pkg/config"
	"github.com/nizar0x1f/termup/pkg/s3storage"
	"github.com/nizar0x1f/termup/pkg/ui"
	"github.com/nizar0x1f/termup/pkg/update"
	"github.com/nizar0x1f/termup/pkg/version"
)

func main() {

	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		showVersion()
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		showHelp()
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "--update" || os.Args[1] == "update") {
		runUpdate()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "relogin" {
		runConfigUI()
		return
	}

	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	filePath := os.Args[1]

	configExists, err := config.Exists()
	if err != nil {
		fmt.Printf("Error checking for config file: %v\n", err)
		os.Exit(1)
	}

	var cfg *config.Config
	if !configExists {

		cfg = runConfigUI()
		if cfg == nil {
			os.Exit(1)
		}
	} else {
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	}

	runUploadUI(cfg, filePath)
}

func runConfigUI() *config.Config {
	model := ui.NewConfigModel()
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running config UI: %v\n", err)
		return nil
	}

	configModel := finalModel.(ui.ConfigModel)
	if !configModel.IsFinished() {
		return nil
	}

	cfg := configModel.GetConfig()
	if cfg == nil {
		return nil
	}

	err = config.Save(cfg)
	if err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return nil
	}

	return cfg
}

func runUploadUI(cfg *config.Config, filePath string) {

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Error accessing file: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewUploadModel(filePath, fileInfo.Size())
	p := tea.NewProgram(model)

	go func() {
		url, err := s3storage.UploadWithProgress(cfg, filePath, func(uploaded int64) {
			p.Send(ui.UploadProgressMsg(uploaded))
		})
		if err != nil {
			p.Send(ui.UploadErrorMsg(err))
		} else {
			p.Send(ui.UploadCompleteMsg(url))
		}
	}()

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running upload UI: %v\n", err)
		os.Exit(1)
	}

	uploadModel := finalModel.(ui.UploadModel)
	if uploadModel.GetError() != nil {
		os.Exit(1)
	}
}

func showVersion() {
	info := version.Get()
	fmt.Println(info.Short())
	fmt.Println("S3 compatible filesharing from terminal")
	fmt.Printf("Built with %s on %s\n", info.GoVersion, info.Platform)
	if info.IsRelease() {
		fmt.Printf("Release: %s (commit: %s)\n", info.Date, info.Commit)
	} else {
		fmt.Printf("Development build (commit: %s)\n", info.Commit)
	}
}

func showHelp() {
	fmt.Println("TermUp - S3 Compatible Filesharing from Terminal")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    upl [OPTIONS] <file-path>")
	fmt.Println("    upl [COMMAND]")
	fmt.Println()
	fmt.Println("ARGS:")
	fmt.Println("    <file-path>    Path to the file to upload")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("    -h, --help       Print help information")
	fmt.Println("    -v, --version    Print version information")
	fmt.Println("        --update     Update to the latest version")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("    relogin          Reconfigure S3 credentials")
	fmt.Println("    update           Update to the latest version")
	fmt.Println("    help             Print this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("    upl document.pdf")
	fmt.Println("    upl photo.jpg")
	fmt.Println("    upl relogin")
	fmt.Println()
	fmt.Println("SUPPORTED PROVIDERS:")
	fmt.Println("    Cloudflare R2, AWS S3, MinIO, DigitalOcean Spaces")
	fmt.Println("    and any other S3-compatible storage service")
}

func showUsage() {
	fmt.Println("Usage: upl <file-path>")
	fmt.Println("       upl relogin")
	fmt.Println("       upl --help")
	fmt.Println("       upl --version")
	fmt.Println("       upl --update")
}

func runUpdate() {
	fmt.Println("Checking for updates...")

	canUpdate, message := update.CanSelfUpdate()
	if !canUpdate {
		fmt.Printf("Self-update not available: %s\n", message)
		fmt.Printf("\nPlease update manually using: %s\n", update.GetUpdateCommand())
		return
	}

	updateInfo, err := update.CheckForUpdates()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		fmt.Println("\nYou can manually update using:")
		fmt.Printf("  %s\n", update.GetUpdateCommand())
		os.Exit(1)
	}

	if !updateInfo.Available {
		fmt.Printf("✅ You're already running the latest version (%s)\n", updateInfo.CurrentVersion)
		return
	}

	fmt.Printf("🚀 Update available!\n")
	fmt.Printf("Current version: %s\n", updateInfo.CurrentVersion)
	fmt.Printf("Latest version:  %s\n", updateInfo.LatestVersion)
	fmt.Printf("\nRelease notes: %s\n\n", updateInfo.ReleaseURL)

	fmt.Print("Do you want to update now? (y/N): ")
	var response string
	_, _ = fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
		fmt.Println("Update cancelled.")
		fmt.Printf("\nTo update later, run: upl --update\n")
		fmt.Printf("Or manually: %s\n", update.GetUpdateCommand())
		return
	}

	fmt.Println("\n⬇️  Updating TermUp...")
	if err := update.PerformSelfUpdate(); err != nil {
		fmt.Printf("❌ Update failed: %v\n", err)
		fmt.Printf("\nPlease update manually using: %s\n", update.GetUpdateCommand())
		os.Exit(1)
	}

	fmt.Println("✅ Update completed successfully!")
	fmt.Println("The new version will be available the next time you run 'upl'")
}
