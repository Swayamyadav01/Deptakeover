# deptakeover 🎯

Ultra-fast supply chain dependency takeover scanner for bug bounty hunting.

Detects unclaimed packages on **npm**, **PyPI**, and **Packagist** registries that can be claimed for supply chain attacks.

## Features

- ✅ **Multi-Ecosystem Support**: Scan npm, Python (PyPI), and PHP (Composer) all at once
- ✅ **Fast**: Built in Go for blazing-fast scanning
- ✅ **Easy to Install**: Single binary, no dependencies
- ✅ **GitHub Integration**: Auto-clone repos from GitHub URLs, usernames, or organizations
- ✅ **Local Scanning**: Scan local projects directly
- ✅ **Detailed Reports**: JSON output with metadata and risk analysis
- ✅ **Bug Bounty Ready**: Identify packages with 404s = takeover targets

## Installation

### Option 1: Download Pre-built Binary

Download the latest binary for your OS:

- **Windows**: `deptakeover-windows-amd64.exe`
- **Linux**: `deptakeover-linux-amd64`
- **macOS (Intel)**: `deptakeover-macos-amd64`
- **macOS (Apple Silicon)**: `deptakeover-macos-arm64`

Extract and add to PATH, then run `deptakeover`.

### Option 2: Build from Source

**Requirements**: Go 1.21+

```bash
# Clone the repo
git clone https://github.com/yourusername/deptakeover.git
cd deptakeover

# Install dependencies
go mod download

# Build for your OS
make build-windows      # Windows
make build-linux        # Linux
make build-macos        # macOS

# Or build for all platforms
make build
```

### Option 3: Direct Installation

```bash
go install github.com/yourusername/deptakeover/cmd/deptakeover@latest
```

## Usage

### Scan a GitHub Repository

```bash
# Scan specific repo
deptakeover --github-repo lodash/lodash --out report.json

# Scan full GitHub URL
deptakeover --github-url https://github.com/django/django --out report.json

# Scan organization (all public repos)
deptakeover --github-org lodash --out report.json

# Scan only npm dependencies
deptakeover --github-repo react/react --ecosystem npm --out report.json

# Scan only Python dependencies
deptakeover --github-repo django/django --ecosystem pypi --out report.json

# Scan only PHP dependencies
deptakeover --github-repo laravel/laravel --ecosystem composer --out report.json
```

### Scan Local Project

```bash
# Scan all ecosystems
deptakeover --local-path ./my-project --out report.json

# Scan specific ecosystem
deptakeover --local-path ./my-project --ecosystem npm,pypi --out report.json
```

## Command Options

| Option | Short | Description | Example |
|--------|-------|-------------|---------|
| `--github-repo` | `-r` | GitHub repo in owner/repo format | `-r lodash/lodash` |
| `--github-url` | `-u` | Full GitHub repo URL | `-u https://github.com/lodash/lodash` |
| `--github-org` | `-o` | GitHub organization name | `-o lodash` |
| `--local-path` | `-p` | Local directory to scan | `-p ./my-project` |
| `--ecosystem` | `-e` | Ecosystems to scan (comma-separated) | `-e npm,pypi,composer` |
| `--out` | `-f` | Output JSON report file | `-f report.json` |

## Output Example

### Console Output

```
🔍 Starting deptakeover supply chain scan...
Scanning: .github_repos/requests

[PyPI] Finding Python dependency files...
[PyPI] Found 3 dependency file(s)
[PyPI] Found 13 unique dependencies
[PyPI] Analyzing packages...

Report saved to: supply_chain_report.json

============================================================
SUPPLY CHAIN TAKEOVER RISK SUMMARY
============================================================
Total dependencies scanned: 13
Total not found (takeover risk): 1

[PYPI]
  Dependencies: 13
  Not found: 1
  ⚠️  Takeover vulnerabilities:
     - fake-package-xyz

============================================================
```

### JSON Report

```json
{
  "repo_path": "...",
  "github_repo": "django/django",
  "ecosystems": {
    "pypi": {
      "total_dependencies": 13,
      "dependencies_by_file": { ... },
      "risk_analysis": {
        "django": {
          "package": "django",
          "exists": true,
          "risk_score": 0,
          "signals": [],
          "metadata": { ... }
        },
        "fake-package-xyz": {
          "package": "fake-package-xyz",
          "exists": false,
          "risk_score": 100,
          "signals": ["not_found_on_pypi"],
          "metadata": {}
        }
      },
      "summary": {
        "not_found_count": 1,
        "not_found_packages": ["fake-package-xyz"],
        "high_risk_count": 0,
        "medium_risk_count": 0
      }
    }
  }
}
```

## How It Works

1. **Repository Discovery**: Clone GitHub repo or scan local directory
2. **Dependency Extraction**: Find dependency files:
   - `package.json` (npm)
   - `requirements.txt`, `setup.py`, `pyproject.toml`, `Pipfile` (Python)
   - `composer.json` (PHP)
3. **Registry Checks**: Query package registries for 404 responses
   - npm: https://registry.npmjs.org/{package}
   - PyPI: https://pypi.org/pypi/{package}/json
   - Packagist: https://packagist.org/packages/{vendor}/{package}.json
4. **Risk Scoring**:
   - 404 (not found) = 100 risk score = **TAKEOVER TARGET** ⚠️
   - Package exists = 0 risk score = Safe

## Examples

### Find Typosquatting Opportunities in Popular Library

```bash
deptakeover --github-repo redux/redux --ecosystem npm --out redux_report.json
cat redux_report.json | jq '.ecosystems.npm.summary.not_found_packages'
```

### Scan Enterprise Monorepo for All Ecosystems

```bash
deptakeover --local-path /path/to/enterprise-repo --ecosystem npm,pypi,composer --out enterprise_report.json
```

### Bulk Scan Organization

```bash
deptakeover --github-org my-org --out org_report.json
```

## Performance

- **Speed**: Scans 100+ dependencies in ~1 minute
- **Concurrency**: Uses goroutines for parallel registry checks
- **Memory**: Minimal memory footprint
- **Network**: Respects registry rate limits

## Bug Bounty Tips

1. **Target Popular Projects**: Scan popular repos like React, Django, Laravel
2. **Look for Typos**: Typosquatting-prone names with misspellings
3. **Check Dependencies**: Include `devDependencies` for additional vectors
4. **Monitor Updates**: New versions = new dependencies = new opportunities
5. **Document Findings**: Always report with evidence (404 response codes)

## Legal & Ethics

- ✅ **Legal**: Only identifying unclaimed packages (no code injection)
- ✅ **Ethical**: Follow responsible disclosure practices
- ⚠️ **Note**: Check bug bounty program terms before claiming packages

## Troubleshooting

### Git not found
If cloning fails, deptakeover automatically falls back to downloading ZIP archives.

### Rate limits
Large scans may hit registry rate limits. Consider adding delays or splitting scans.

### Permission denied (local path)
Ensure the scanning user has read access to the directory.

## Contributing

```bash
# Fork and clone
git clone https://github.com/yourusername/deptakeover.git
cd deptakeover

# Create feature branch
git checkout -b feature/your-feature

# Make changes and test
go test ./...

# Build and test
make build

# Submit PR
```

## License

MIT License - See LICENSE file

## Author

Created for bug bounty hunters and security researchers.

## Support

- 🐛 Report bugs: [GitHub Issues]
- 💡 Suggest features: [GitHub Discussions]
- 📧 Email: contact@example.com

---

**Happy hunting! 🎯**
