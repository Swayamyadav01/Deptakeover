# Contributing

Pull requests welcome. Here's how to contribute:

## Setup

```bash
git clone https://github.com/Swayamyadav01/Deptakeover.git
cd Deptakeover
go mod download
```

## Making Changes

- Fork the repo
- Make your changes
- Test them (`go test ./...`)
- Submit a PR with a clear description

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing  
- Add tests for new features
- Keep it simple

# Build and test manually
go build -o deptakeover ./cmd/deptakeover
./deptakeover npm lodash/lodash
```

### 5. Commit & Push
```bash
git add .
git commit -m "Add feature: description of changes"
git push origin feature/your-feature-name
```

### 6. Create Pull Request
- Open a PR against the `main` branch
- Provide clear description of changes
- Reference any related issues

## 📋 Contribution Guidelines

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small

### Testing
- Add unit tests for new functionality
- Test with real repositories when possible
- Ensure cross-platform compatibility

### Documentation
- Update README.md for new features
- Add inline code documentation
- Include usage examples

## 🐛 Bug Reports

When reporting bugs, please include:
- Operating system and version
- Go version
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs

**Template:**
```markdown
**Environment:**
- OS: Windows 11 / Ubuntu 22.04 / macOS 13
- Go version: 1.21.0
- DepTakeover version: v1.0.0

**Bug Description:**
Clear description of the issue

**Steps to Reproduce:**
1. Run command X
2. Observe behavior Y
3. Expected Z

**Error Output:**
```
paste error messages here
```
```

## 💡 Feature Requests

We welcome new feature ideas! Please open an issue with:
- Clear description of the feature
- Use case or problem it solves
- Proposed implementation (if you have ideas)
- Willingness to implement it yourself

### Priority Areas
- New package registry support (RubyGems, NuGet, etc.)
- Performance optimizations
- Enhanced reporting features
- Better error handling
- CI/CD integrations

## 🔒 Security Issues

**DO NOT** open public issues for security vulnerabilities. Instead:
1. Email security concerns to: [security@example.com]
2. Follow our [Security Policy](SECURITY.md)
3. Allow time for assessment and fix before disclosure

## 📖 Development Setup

### Prerequisites
- Go 1.21 or later
- Git
- Text editor or IDE

### Local Development
```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build -o deptakeover ./cmd/deptakeover

# Cross-compile for all platforms
make build-all  # or use build.bat / build.sh
```

### Project Structure
```
deptakeover/
├── cmd/deptakeover/      # Main application entry point
├── internal/             # Internal packages
│   ├── scanner/         # Dependency file parsers
│   ├── registry/        # Registry API clients
│   └── github/          # GitHub integration
├── docs/                # Documentation
├── build/               # Build scripts and outputs
└── tests/               # Test data and integration tests
```

### Adding New Registry Support

To add support for a new package registry:

1. **Create registry client** in `internal/registry/`
   ```go
   // internal/registry/newregistry.go
   func CheckNewRegistryPackageRisk(pkg string) PackageInfo {
       // Implementation
   }
   ```

2. **Add scanner support** in `internal/scanner/`
   ```go
   // internal/scanner/newlang.go
   func ExtractNewLangDependencies(repoPath string) map[string][]string {
       // Implementation
   }
   ```

3. **Update main CLI** in `cmd/deptakeover/main.go`
   - Add ecosystem mapping
   - Add scanning logic
   - Update help text

4. **Add tests** for new functionality

5. **Update documentation**

## 🏆 Recognition

Contributors will be:
- Listed in the project contributors
- Credited in release notes
- Mentioned in security advisories (if applicable)

## 📜 Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## 📞 Questions?

- Open a GitHub issue for bugs or feature requests
- Join our discussions in GitHub Discussions
- Tag @maintainers for urgent issues

---

Thank you for helping make DepTakeover better! 🚀