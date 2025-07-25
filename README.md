# TermUp - S3 Compatible Filesharing from Terminal

A modern, fast, and beautiful command-line utility for S3-compatible filesharing with real-time progress tracking and a stunning terminal UI. Upload files to any S3-compatible storage service including Cloudflare R2, AWS S3, MinIO, and more.

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)

## Features

- **S3 Compatible** - Works with any S3-compatible storage (R2, AWS S3, MinIO, DigitalOcean Spaces, etc.)
- **Beautiful Terminal UI** - Powered by Charm's Bubble Tea framework
- **Real-time Progress** - Live transfer speed, ETA, and progress visualization
- **Fast Uploads** - Optimized transfers with parallel processing
- **Secure Configuration** - Encrypted credential storage
- **Custom Domains** - Support for custom public URLs and CDNs
- **Cross-platform** - Works on macOS, Linux, and Windows
- **Easy Reconfiguration** - Simple credential management
- **Multiple Providers** - Switch between different S3 providers easily

## Demo

```
S3 Compatible File Upload

Uploading: document.pdf

⣻  ████████████████████░░░░░░░░░░░░░░░  54%

2.1 MB / 4.0 MB (52.5%) 1.2 MB/s ETA: 00:02
Elapsed: 00:01

✓ Upload successful!
URL: https://your-s3-domain.com/document.pdf
```

## Supported S3 Providers

TermUp works with any S3-compatible storage service:

- **Cloudflare R2** - Fast, global object storage
- **AWS S3** - Amazon's original object storage service
- **MinIO** - High-performance, self-hosted object storage
- **DigitalOcean Spaces** - Simple, scalable object storage
- **Backblaze B2** - Affordable cloud storage
- **Wasabi** - Hot cloud storage
- **Linode Object Storage** - Simple, scalable storage
- **Any S3-compatible service** - If it speaks S3 API, it works!

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/nizar0x1f/termup/cmd/upl@latest
```



### From Source

```bash
git clone https://github.com/nizar0x1f/termup.git
cd termup
go install -v ./cmd/upl
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/nizar0x1f/termup/releases).

## Quick Start

### First Time Setup

Run the tool for the first time to configure your S3-compatible storage credentials:

```bash
upl myfile.txt
```

You'll be prompted to enter:
- **Access Key ID** - Your S3 access key ID
- **Secret Access Key** - Your S3 secret access key
- **Bucket Name** - The bucket to upload to
- **Endpoint** - Your S3-compatible endpoint URL
- **Public URL** - Public URL for accessing files (optional)

**Tip:** You can paste values using Ctrl+V, Cmd+V, or right-click. The interface automatically handles bracketed paste and removes unwanted characters.

### Basic Usage

```bash
# Upload a file
upl document.pdf

# Upload an image
upl photo.jpg

# Upload any file type
upl archive.zip

# Check version
upl --version

# Get help
upl --help

# Update to latest version
upl --update
```

## Configuration

### Reconfigure Credentials

```bash
upl relogin
```

This will prompt you to re-enter all configuration details.

### Configuration File

TermUp stores configuration in `~/.termup.json`:

```json
{
  "access_key_id": "your-access-key",
  "secret_access_key": "your-secret-key",
  "bucket": "your-bucket-name",
  "endpoint": "https://your-s3-endpoint.com",
  "public_url": "https://your-custom-domain.com/"
}
```

### S3 Provider Examples

**Cloudflare R2:**
```json
{
  "endpoint": "https://your-account-id.r2.cloudflarestorage.com",
  "public_url": "https://pub-xyz.r2.dev/"
}
```

**AWS S3:**
```json
{
  "endpoint": "https://s3.us-east-1.amazonaws.com",
  "public_url": "https://your-bucket.s3.amazonaws.com/"
}
```

**MinIO:**
```json
{
  "endpoint": "https://your-minio-server.com",
  "public_url": "https://your-minio-server.com/your-bucket/"
}
```

**DigitalOcean Spaces:**
```json
{
  "endpoint": "https://nyc3.digitaloceanspaces.com",
  "public_url": "https://your-space.nyc3.cdn.digitaloceanspaces.com/"
}
```

**Backblaze B2:**
```json
{
  "endpoint": "https://s3.us-west-000.backblazeb2.com",
  "public_url": "https://f000.backblazeb2.com/file/your-bucket/"
}
```

**Wasabi:**
```json
{
  "endpoint": "https://s3.wasabisys.com",
  "public_url": "https://s3.wasabisys.com/your-bucket/"
}
```

**Linode Object Storage:**
```json
{
  "endpoint": "https://us-east-1.linodeobjects.com",
  "public_url": "https://your-bucket.us-east-1.linodeobjects.com/"
}
```

## UI Features

### Beautiful Configuration Interface

- **Step-by-step setup** with clear prompts
- **Secure password input** with masked characters
- **Default value suggestions** for common settings
- **Colorful, modern terminal interface**

### Advanced Progress Display

- **Real-time transfer speed** (B/s, KB/s, MB/s, GB/s)
- **Animated progress bar** with gradient colors
- **ETA calculation** based on current speed
- **Elapsed time tracking**
- **File size and percentage display**
- **Spinning activity indicator**

### Error Handling

- **Clear error messages** with helpful suggestions
- **Network error recovery** with retry logic
- **Invalid credential detection**
- **File access error reporting**

## Updates

### Checking for Updates

```bash
# Check current version
upl --version
```

### Updating TermUp

**Built-in Update (Recommended):**
```bash
upl --update
```

**Manual Update:**
```bash
go install github.com/nizar0x1f/termup/cmd/upl@latest
```

**Update to Specific Version:**
```bash
go install github.com/nizar0x1f/termup/cmd/upl@v1.2.0
```

### Version Information

TermUp follows [Semantic Versioning](https://semver.org/). Version numbers follow the format `MAJOR.MINOR.PATCH`:

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in a backwards compatible manner
- **PATCH**: Backwards compatible bug fixes

## Advanced Usage

### Environment Variables

You can override configuration using environment variables:

```bash
export S3_ACCESS_KEY_ID="your-key"
export S3_SECRET_ACCESS_KEY="your-secret"
export S3_BUCKET="your-bucket"
export S3_ENDPOINT="your-endpoint"
export S3_PUBLIC_URL="your-public-url"

upl myfile.txt
```

### Batch Operations

```bash
# Upload multiple files
for file in *.jpg; do
  upl "$file"
done

# Upload with custom naming
upl important-document.pdf
# Returns: https://your-domain.com/important-document.pdf
```

### Integration Examples

```bash
# Use in scripts
URL=$(upl screenshot.png | grep "URL:" | cut -d' ' -f2)
echo "File uploaded to: $URL"

# Pipe output
echo "Uploading backup..." && upl backup.tar.gz

# Conditional upload
if upl config.json; then
  echo "Configuration uploaded successfully"
fi
```

## Troubleshooting

### Common Issues

#### "Access Denied" Error
```
Error uploading file: Access Denied
```
**Solutions:**
- Verify your Access Key ID and Secret Access Key
- Check that your API token has write permissions for the bucket
- Ensure the bucket name is correct
- Confirm the endpoint URL matches your account

#### "Bucket Not Found" Error
```
Error uploading file: NoSuchBucket
```
**Solutions:**
- Verify the bucket exists in your R2 dashboard
- Check the bucket name spelling
- Ensure you have access to the bucket

#### "Network Connection" Issues
```
Error uploading file: connection timeout
```
**Solutions:**
- Check your internet connection
- Verify the endpoint URL is correct
- Try again (the tool has built-in retry logic)

#### Configuration Issues
```bash
# Reset configuration
rm ~/.termup.json
upl relogin
```

### Debug Mode

For detailed error information, check the upload process:

```bash
# Verbose output
go run ./cmd/upl -v myfile.txt

# Check configuration
cat ~/.termup.json
```

## Development

### Prerequisites

- Go 1.22 or later
- Git

### Releasing

TermUp uses automated releases triggered by Git tags. To create a new release:

**Using the release script (recommended):**
```bash
./scripts/release.sh v1.2.0
```

**Manual process:**
```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# Run tests
CI=true go test ./...

# Create and push tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

The release process automatically:
- Builds binaries for all platforms (Linux, macOS, Windows)
- Creates GitHub release with assets and changelog
- Updates Homebrew tap for easy installation
- Generates checksums for verification

### Building from Source

```bash
# Clone the repository
git clone https://github.com/nizar0x1f/termup.git
cd termup

# Install dependencies
go mod tidy

# Build the binary
go build -o upl ./cmd/upl

# Run tests
go test ./...

# Install locally
go install -v ./cmd/upl
```

### Project Structure

```
termup/
├── cmd/upl/           # Main application
│   ├── main.go        # Entry point
│   └── main_test.go   # Main tests
├── pkg/
│   ├── config/        # Configuration management
│   ├── s3storage/     # S3-compatible upload logic
│   └── ui/           # Terminal UI components
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
└── README.md         # This file
```

### Dependencies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - UI components
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling library
- **[AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)** - S3-compatible API client
- **[Progress Bar](https://github.com/cheggaaa/pb)** - Fallback progress display

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/config
go test ./pkg/r2
go test ./pkg/ui
```

## Contributing

We welcome contributions! Here's how you can help:

### Reporting Issues

1. Check existing [issues](https://github.com/nizar0x1f/termup/issues)
2. Create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - System information (OS, Go version)

### Submitting Changes

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**:
   - Follow Go conventions
   - Add tests for new functionality
   - Update documentation
4. **Test your changes**:
   ```bash
   go test ./...
   go build ./cmd/upl
   ```
5. **Commit your changes**:
   ```bash
   git commit -m "Add amazing feature"
   ```
6. **Push to your fork**:
   ```bash
   git push origin feature/amazing-feature
   ```
7. **Create a Pull Request**

### Development Guidelines

- **Code Style**: Follow `gofmt` and `golint` standards
- **Testing**: Add tests for new features
- **Documentation**: Update README and code comments
- **Commits**: Use clear, descriptive commit messages
- **Dependencies**: Minimize external dependencies

### Feature Requests

We're always looking for ways to improve TermUp! Consider contributing:

- **New cloud storage providers** (AWS S3, Google Cloud Storage, etc.)
- **Additional file formats** and validation
- **Performance optimizations**
- **UI/UX improvements**
- **Cross-platform enhancements**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 TermUp Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```


## Links

- **[GitHub Repository](https://github.com/nizar0x1f/termup)**
- **[Issue Tracker](https://github.com/nizar0x1f/termup/issues)**
- **[Releases](https://github.com/nizar0x1f/termup/releases)**

### S3 Provider Documentation
- **[AWS S3 Documentation](https://docs.aws.amazon.com/s3/)**
- **[Cloudflare R2 Documentation](https://developers.cloudflare.com/r2/)**
- **[MinIO Documentation](https://min.io/docs/)**
- **[DigitalOcean Spaces Documentation](https://docs.digitalocean.com/products/spaces/)**

### Development
- **[Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)**
- **[AWS SDK for Go](https://aws.github.io/aws-sdk-go-v2/)**

## Stats

![GitHub stars](https://img.shields.io/github/stars/nizar0x1f/termup?style=social)
![GitHub forks](https://img.shields.io/github/forks/nizar0x1f/termup?style=social)
![GitHub issues](https://img.shields.io/github/issues/nizar0x1f/termup)
![GitHub pull requests](https://img.shields.io/github/issues-pr/nizar0x1f/termup)

---

<div align="center">
 <sub>S3 compatible filesharing from terminal - Built with Go and Bubble Tea</sub>
</div>
