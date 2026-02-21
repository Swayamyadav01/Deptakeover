# 🔍 DepTakeover - Supply Chain Takeover Scanner

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/yourusername/deptakeover/releases)
[![Bug Bounty](https://img.shields.io/badge/Built%20for-Bug%20Bounty-red?style=for-the-badge)](https://github.com/yourusername/deptakeover)

> ⚡ **Lightning-fast supply chain vulnerability scanner designed for bug bounty hunters**

Hunt for **package takeover vulnerabilities** across npm, PyPI, and Composer registries. Scan individual repositories or entire organizations to discover unclaimed packages that could lead to supply chain attacks.

## 🎯 What is Package Takeover?

When a project depends on a package that **no longer exists** on the registry (npm, PyPI, Packagist), an attacker can claim that package name and potentially compromise all projects that depend on it. This scanner identifies these vulnerable dependencies automatically.

## ✨ Features

- 🚀 **Blazing Fast** - Written in Go for maximum performance
- 🌐 **Multi-Registry Support** - npm, PyPI (Python), Composer (PHP)
- 🏢 **Organization Scanning** - Scan entire GitHub organizations
- 📊 **Comprehensive Reports** - Detailed JSON reports with vulnerability analysis
- 🎯 **Bug Bounty Optimized** - Built specifically for security researchers
- 💻 **Cross-Platform** - Windows, Linux, macOS support
- ⚡ **Zero Dependencies** - Single binary, no runtime requirements

## 🚀 Quick Start

### Installation

**Option 1: Download Pre-built Binary**
```bash
# Download from releases page
curl -L https://github.com/yourusername/deptakeover/releases/latest/download/deptakeover-linux-amd64 -o deptakeover
chmod +x deptakeover
```

**Option 2: Build from Source**
```bash
git clone https://github.com/yourusername/deptakeover.git
cd deptakeover
go build -o deptakeover ./cmd/deptakeover
```

### Basic Usage

**Single Repository Scanning:**
```bash
# Scan npm dependencies
deptakeover npm lodash/lodash

# Scan Python packages
deptakeover pypi django/django

# Scan PHP packages  
deptakeover composer laravel/laravel

# Shorthand aliases
deptakeover py requests/requests    # Python
deptakeover php symfony/symfony     # PHP
```

**Organization-Wide Scanning:**
```bash
# Scan all ecosystems across entire org
deptakeover org microsoft

# Ecosystem-specific org scans
deptakeover org-npm facebook        # npm only
deptakeover org-pypi google         # Python only
deptakeover org-composer symfony    # PHP only
```

## 📊 Example Output

```bash
$ deptakeover npm lodash/lodash

🔍 Scanning [npm]...
📦 Found 27 packages
✅ Report: npm_report.json
──────────────────────────────────────────────────
📊 Dependencies: 27
⚠️  Takeover targets: 0
──────────────────────────────────────────────────
```

**With Vulnerabilities Found:**
```bash
$ deptakeover composer vulnerable-project/repo

🔍 Scanning [composer]...  
📦 Found 15 packages
✅ Report: composer_report.json
──────────────────────────────────────────────────
📊 Dependencies: 15
⚠️  Takeover targets: 3

🚨 [COMPOSER] 3 NOT FOUND:
  • abandoned-package/helper
  • old-vendor/legacy-lib
  • missing-dep/validator
──────────────────────────────────────────────────
```

## 🏢 Organization Scanning

Perfect for discovering vulnerabilities across large organizations:

```bash
# Scan Microsoft's repositories
deptakeover org microsoft

# Results show:
# - Total repositories scanned: 2,847
# - Total vulnerabilities: 23
# - Top vulnerable packages
# - Most vulnerable repositories
# - Frequency analysis
```

## 📁 Project Structure

```
deptakeover/
├── cmd/deptakeover/          # Main application
├── internal/
│   ├── scanner/              # Dependency file parsers
│   │   ├── npm.go           # package.json parser
│   │   ├── python.go        # requirements.txt, setup.py, etc.
│   │   └── php.go           # composer.json parser
│   ├── registry/             # Registry API clients
│   │   ├── npm.go           # npm registry checks
│   │   ├── pypi.go          # PyPI registry checks  
│   │   └── packagist.go     # Packagist registry checks
│   └── github/              # GitHub integration
│       └── handler.go       # Repository cloning/downloading
├── build/                   # Build outputs
├── docs/                    # Documentation
└── README.md
```

## 🔧 Advanced Usage

### CI/CD Integration

```yaml
# .github/workflows/security-scan.yml
name: Supply Chain Security Scan
on: [push, pull_request]
jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Download DepTakeover
        run: |
          curl -L https://github.com/yourusername/deptakeover/releases/latest/download/deptakeover-linux-amd64 -o deptakeover
          chmod +x deptakeover
      - name: Scan Dependencies
        run: ./deptakeover npm ${{ github.repository }}
```

### Custom Rate Limiting

For large organizations, the tool automatically handles rate limiting:
- GitHub API: 500ms between repos
- Registry APIs: Parallel requests with backoff
- Configurable timeouts for large repositories

### Report Analysis

JSON reports include:
- **Vulnerability Details**: Package names, risk scores, registry status
- **Repository Metadata**: Stars, language, size
- **Dependency Analysis**: File locations, version requirements
- **Risk Assessment**: High/medium/low risk categorization

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
git clone https://github.com/yourusername/deptakeover.git
cd deptakeover
go mod download
go build -o deptakeover./cmd/deptakeover
```

### Running Tests

```bash
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

Found a security issue? Please see our [Security Policy](SECURITY.md) for responsible disclosure.

## 🏆 Bug Bounty Tips

1. **Start with org-npm scans** - Faster for JavaScript-heavy organizations
2. **Focus on popular repositories** - Higher impact vulnerabilities
3. **Check dependency frequency** - Packages used across multiple repos
4. **Verify manual takeover** - Always confirm packages are truly unclaimed
5. **Document impact** - Show potential reach of supply chain attack

## 📈 Roadmap

- [ ] Support for additional registries (RubyGems, NuGet, etc.)
- [ ] Real-time monitoring mode
- [ ] Web dashboard for results visualization
- [ ] Integration with popular bug bounty platforms
- [ ] Automated proof-of-concept generation

## 🙏 Acknowledgments

- Built for the bug bounty and security research community
- Inspired by dependency confusion research by Alex Birsan
- Thanks to all contributors and security researchers

---

⭐ **Star this repo if it helped you find vulnerabilities!**

**Disclaimer**: This tool is for authorized security research only. Always follow responsible disclosure practices and respect bug bounty program terms.