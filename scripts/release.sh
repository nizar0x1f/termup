#!/bin/bash

# TermUp Release Script
# Usage: ./scripts/release.sh [version]
# Example: ./scripts/release.sh v1.2.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version is provided
if [ $# -eq 0 ]; then
    print_error "Version is required"
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.0"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format. Use semantic versioning (e.g., v1.2.0)"
    exit 1
fi

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "You're not on the main branch (current: $CURRENT_BRANCH)"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Release cancelled"
        exit 1
    fi
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    print_error "Working directory is not clean. Please commit or stash your changes."
    git status --short
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    print_error "Tag $VERSION already exists"
    exit 1
fi

# Fetch latest changes
print_status "Fetching latest changes..."
git fetch origin

# Check if local main is up to date with remote
LOCAL=$(git rev-parse main)
REMOTE=$(git rev-parse origin/main)
if [ "$LOCAL" != "$REMOTE" ]; then
    print_error "Local main branch is not up to date with origin/main"
    print_status "Run: git pull origin main"
    exit 1
fi

# Run tests
print_status "Running tests..."
if ! CI=true go test ./...; then
    print_error "Tests failed. Please fix them before releasing."
    exit 1
fi

print_success "All tests passed!"

# Build to ensure everything compiles
print_status "Building project..."
if ! make build; then
    print_error "Build failed. Please fix compilation errors before releasing."
    exit 1
fi

print_success "Build successful!"

# Show recent commits for release notes
print_status "Recent commits since last tag:"
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LAST_TAG" ]; then
    git log --oneline "$LAST_TAG"..HEAD
else
    git log --oneline -10
fi

echo
print_status "Creating release $VERSION"
read -p "Continue with release? (y/N): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_error "Release cancelled"
    exit 1
fi

# Create and push tag
print_status "Creating tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

print_status "Pushing tag to origin..."
git push origin "$VERSION"

print_success "Release $VERSION has been triggered!"
print_status "Check the GitHub Actions tab for build progress:"
print_status "https://github.com/nizar0x1f/termup/actions"

print_success "Release will be available at:"
print_status "https://github.com/nizar0x1f/termup/releases/tag/$VERSION"
