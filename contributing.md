# 🤝 Contributing to GORR

> _Thank you for your interest in contributing to GORR! This document provides guidelines for contributing to the project._

<div align="center">

[![Contributors Welcome](https://img.shields.io/badge/Contributors-Welcome-brightgreen?style=flat-square)](CONTRIBUTING.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-blue?style=flat-square)](https://github.com/GlazKrovi/go-re-releaser/pulls)

</div>

## 🎯 Quick Start

Want to contribute but don't know where to start? Check out our [Good First Issues](https://github.com/GlazKrovi/go-re-releaser/labels/good%20first%20issue) label!

## 🛠️ Development Setup

### 📋 Prerequisites

- 🐹 **Go 1.21 or later**
- 📁 **Git**
- 📦 **GoReleaser** (for testing)

### 🔨 Building the Project

```bash
# 📥 Clone the repository
git clone <repository-url>
cd go-re-releaser

# 🔨 Build the binary
go build -o gorr ./cmd/gorr

# ✅ Test the build
./gorr --version
```

### 🧪 Testing Your Changes

```bash
# 🧪 Test the release command (snapshot mode)
./gorr release patch --snapshot

# 🔧 Test direct GoReleaser commands
./gorr check
./gorr init
```

## 🔄 Development Workflow

### 1. 🍴 Fork and Clone

```bash
# 🍴 Fork the repository on GitHub
# 📥 Clone your fork
git clone https://github.com/GlazKrovi/go-re-releaser.git
cd go-re-releaser
```

### 2. 🌿 Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 3. ✏️ Make Your Changes

- 📝 Follow Go coding standards
- 🧪 Add tests for new functionality
- 📚 Update documentation as needed
- ✅ Ensure all existing tests pass

### 4. 🧪 Test Your Changes

```bash
# 🔨 Build and test
go build -o gorr ./cmd/gorr
./gorr release patch --snapshot

# 🧪 Run any additional tests
go test ./...
```

### 5. 💾 Commit Your Changes

```bash
git add .
git commit -m "feat: add new feature description"
```

### 6. 🚀 Push and Create Pull Request

```bash
git push origin feature/your-feature-name
# 📝 Create PR on GitHub
```

## 📝 Code Style Guidelines

### 🐹 Go Code Style

- 🎨 Use `gofmt` to format code
- 📛 Follow standard Go naming conventions
- 💬 Add comments for exported functions
- 🎯 Keep functions focused and small

### 🚨 Error Handling

- ⚠️ Always handle errors explicitly
- 💬 Provide meaningful error messages
- 🔗 Use `fmt.Errorf` for error wrapping

### 💡 Example

```go
func exampleFunction() error {
    cmd := exec.Command("git", "status")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("git status failed: %v", err)
    }
    return nil
}
```

## 🧪 Testing Guidelines

### 🔬 Unit Tests

- 🎯 Test individual functions
- 🎭 Mock external dependencies when possible
- 📊 Aim for high test coverage

### 🔗 Integration Tests

- 🔄 Test the full workflow
- 🧪 Use snapshot mode for testing
- ✅ Verify Git operations work correctly

### 💡 Example Test

```go
func TestGetNextVersion(t *testing.T) {
    tests := []struct {
        current    string
        releaseType string
        expected   string
    }{
        {"v1.0.0", "patch", "v1.0.1"},
        {"v1.0.0", "minor", "v1.1.0"},
        {"v1.0.0", "major", "v2.0.0"},
    }

    for _, test := range tests {
        result := getNextVersion(test.current, test.releaseType)
        if result != test.expected {
            t.Errorf("Expected %s, got %s", test.expected, result)
        }
    }
}
```

## 📝 Pull Request Guidelines

### ✅ Before Submitting

- [ ] 📝 Code follows Go style guidelines
- [ ] ✅ All tests pass
- [ ] 📚 Documentation is updated
- [ ] 💬 Commit messages are clear and descriptive

### 📋 PR Description

Include:

- 🔄 What changes were made
- 🤔 Why the changes were necessary
- 🧪 How to test the changes
- ⚠️ Any breaking changes

## 🐛 Issue Reporting

When reporting issues:

1. 🔍 **Check existing issues** first
2. 📝 **Provide clear description** of the problem
3. 🔄 **Include steps to reproduce**
4. 💻 **Specify environment details** (OS, Go version, etc.)
5. 📋 **Include relevant logs** or error messages

## 💡 Feature Requests

For new features:

1. 🔍 **Check if feature already exists**
2. 📝 **Describe the use case** clearly
3. 🤔 **Explain why it would be beneficial**
4. ⚙️ **Consider implementation complexity**

## 👀 Code Review Process

### 👨‍💻 For Contributors

- ⚡ Respond to review feedback promptly
- 🔄 Make requested changes clearly
- ❓ Ask questions if feedback is unclear
- 🧪 Test changes after addressing feedback

### 👨‍🔬 For Reviewers

- 💬 Provide constructive feedback
- 🎯 Focus on code quality and maintainability
- 🤝 Be respectful and helpful
- ✅ Approve when ready

## 🚀 Release Process

1. 📈 **Version bumping** is handled automatically
2. 📋 **Changelog** should be updated
3. ✅ **Tests** must pass before release
4. 📚 **Documentation** should be current

---

<div align="center">

**🎉 Thank you for contributing to GORR! 🎉**

[⭐ Star this repo](https://github.com/GlazKrovi/go-re-releaser) • [🐛 Report Bug](https://github.com/GlazKrovi/go-re-releaser/issues) • [💡 Request Feature](https://github.com/GlazKrovi/go-re-releaser/issues)

</div>
