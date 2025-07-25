# Contributing to TermUp

Thank you for your interest in contributing to TermUp - S3 compatible filesharing from terminal! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.24 or later


### Setting Up Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/your-username/termup.git
   cd termup
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/nizar0x1f/termup.git
   ```
4. **Install dependencies**:
   ```bash
   go mod tidy
   ```
5. **Build and test**:
   ```bash
   go build ./cmd/upl
   go test ./...
   ```

## Development Workflow

### Before You Start

1. **Check existing issues** to avoid duplicate work
2. **Create an issue** for new features or bugs
3. **Discuss your approach** in the issue comments

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. **Make your changes** following our coding standards
3. **Write tests** for new functionality
4. **Update documentation** as needed
5. **Test your changes**:
   ```bash
   go test ./...
   go build ./cmd/upl
   ./upl 
   ```

### Submitting Changes

1. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add amazing new feature"
   ```
2. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
3. **Create a Pull Request** on GitHub

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(ui): add real-time progress display
fix(r2): handle connection timeout errors
docs: update installation instructions
```

### Testing

- Write unit tests for new functions
- Test error conditions
- Ensure tests pass before submitting
- Aim for good test coverage

### Documentation

- Update README.md for new features
- Add inline code comments
- Update help text and usage examples
- Include examples in documentation

## Reporting Issues

### Bug Reports

Include:
- **Clear description** of the problem
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **System information** (OS, Go version, etc.)
- **Error messages** or logs
- **Screenshots** if applicable 

### Feature Requests

Include:
- **Clear description** of the feature
- **Use case** and motivation
- **Proposed implementation** (if you have ideas)
- **Alternatives considered**

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/config
```

### Test Structure

- Place tests in `*_test.go` files
- Use table-driven tests when appropriate
- Mock external dependencies
- Test both success and error cases

## Documentation

### Code Documentation

- Document all exported functions and types
- Use clear, concise comments
- Include examples in documentation comments
- Follow Go documentation conventions

### User Documentation

- Update README.md for user-facing changes
- Include usage examples
- Document configuration options
- Explain error messages and solutions

## Areas for Contribution

We welcome contributions in these areas:

### High Priority
- Bug fixes and stability improvements
- Performance optimizations
- Cross-platform compatibility
- Error handling improvements

### Medium Priority
- New cloud storage providers
- Additional file validation
- UI/UX enhancements
- Configuration improvements

### Low Priority
- Code refactoring
- Documentation improvements
- Test coverage increases
- Development tooling

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Code Review**: We provide feedback on all pull requests

## Pull Request Checklist

Before submitting a pull request, ensure:

- [ ] Code follows Go conventions and passes `go fmt`
- [ ] All tests pass (`go test ./...`)
- [ ] New functionality includes tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format
- [ ] PR description explains the changes
- [ ] No breaking changes (or clearly documented)

## Recognition

Contributors will be:
- Listed in the project's contributors
- Mentioned in release notes for significant contributions
- Invited to join the maintainer team for sustained contributions

Thank you for contributing to TermUp!
