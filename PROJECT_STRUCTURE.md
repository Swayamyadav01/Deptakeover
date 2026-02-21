# deptakeover - Project Structure

## Directory Layout

```
deptakeover/
├── cmd/
│   └── deptakeover/
│       └── main.go                 # Main CLI entrypoint (Cobra framework)
│
├── internal/
│   ├── github/
│   │   └── handler.go              # GitHub repo cloning & zip fallback
│   │
│   ├── registry/
│   │   ├── npm.go                  # npm registry analyzer
│   │   ├── pypi.go                 # PyPI registry analyzer
│   │   └── packagist.go            # Packagist (PHP) registry analyzer
│   │
│   └── scanner/
│       ├── npm.go                  # package.json parser
│       ├── python.go               # Python dependency file parser
│       └── php.go                  # composer.json parser
│
├── go.mod                          # Go module definition
├── go.sum                          # Dependency lock file
│
├── Makefile                        # Build targets for Unix/Linux
├── build.sh                        # Build script for Unix/Linux/macOS
├── build.bat                       # Build script for Windows
│
├── SETUP_GO.md                     # Setup & installation guide
├── README_GO.md                    # Complete documentation
└── README.md                       # Original project docs (Python)

# Python Version (legacy)
├── src/
│   ├── github_handler.py
│   ├── scanner.py
│   ├── npm_analyzer.py
│   ├── python_scanner.py
│   ├── pypi_analyzer.py
│   ├── php_scanner.py
│   └── packagist_analyzer.py
│
├── cli/
│   └── main.py
│
├── tests/
│   └── test_fresh.py
│
└── requirements.txt
```

## Comparison: Python vs Go

| Feature | Python | Go |
|---------|--------|-----|
| **Speed** | ~2-3 sec per dependency | <100ms per dependency |
| **Binary Size** | ~50MB (with interpreter) | ~8-12MB (single binary) |
| **Installation** | pip install + venv | Download & run |
| **Dependencies** | Multiple packages | None (statically linked) |
| **Startup Time** | 1-2 seconds | <100ms |
| **Total Scan Time (100 deps)** | ~5-10 minutes | ~1-2 minutes |
| **Memory Usage** | 150-200MB | 20-50MB |
| **Distribution** | Source code | Single executable |
| **Cross-platform** | Requires Python runtime | Native for each OS |

## Key Go Modules Used

- **cli/cobra**: Command-line interface framework
- **go-git/go-git**: Git repository cloning
- **net/http**: HTTP requests to registries
- **encoding/json**: JSON parsing & generation
- **regexp**: Pattern matching for dependency parsing

## Build Targets

### Supported Platforms

- Windows (64-bit)
- Linux (64-bit)
- macOS (Intel 64-bit)
- macOS (Apple Silicon ARM64)

### Build Commands

```bash
# Windows
make build-windows

# Linux
make build-linux

# macOS
make build-macos

# All platforms
make build

# Cleanup
make clean
```

## Performance Characteristics

### Scanning Time (Approximate)

- **Small project** (10-20 deps): 5-10 seconds
- **Medium project** (50-100 deps): 1-2 minutes
- **Large project** (200+ deps): 3-5 minutes

**Factors affecting speed:**
- Network latency to registries
- Registry response times
- Timeout settings (5 seconds per package)

### Memory Usage

- Base: ~10MB
- Per 100 dependencies: +10-20MB
- Typical scan: 20-50MB total

## Registry API Endpoints

### npm (JavaScript)
```
https://registry.npmjs.org/{package}
```
- Returns 404 if package not found
- Returns JSON metadata if exists

### PyPI (Python)
```
https://pypi.org/pypi/{package}/json
```
- Returns 404 if package not found
- Returns package metadata if exists

### Packagist (PHP)
```
https://packagist.org/packages/{vendor}/{package}.json
```
- Returns 404 if package not found
- Returns package metadata if exists

## Code Structure

### Registry Package (internal/registry/)

Each registry has functions:
- `Check{Registry}PackageRisk()` - Single package check
- `Analyze{Registry}DependencyRisks()` - Batch analysis

Result structure:
```go
type SPkg struct {
    Package   string                 // Package name
    Exists    bool                   // Found on registry
    RiskScore int                    // 0-100
    Signals   []string               // Risk indicators
    Metadata  map[string]interface{} // Package info
}
```

### Scanner Package (internal/scanner/)

Each scanner has functions:
- `Find{Ecosystem}Files()` - Locate dependency files
- `Extract{Ecosystem}Dependencies()` - Parse dependencies
- `ExtractAll{Ecosystem}Dependencies()` - Multiple files

Result structure:
```go
map[filepath]map[packagename]version
```

### GitHub Package (internal/github/)

- `GetRepoPath()` - Resolve input to repo path
- `CloneGitHubRepo()` - Clone or zip download
- `downloadZip()` - Fallback zip extraction

### Main CLI (cmd/deptakeover/)

- Cobra command setup
- Report generation
- JSON serialization
- Console output formatting

## Dependency Resolution

The tool discovers dependencies from:

### npm
- `package.json` → `dependencies` and `devDependencies`
- Skips `@scoped` packages

### Python
- `requirements.txt` → Each line = package
- `setup.py` → `install_requires` parameter
- `pyproject.toml` → `[project] dependencies`
- `Pipfile` → `[packages]` section

### PHP
- `composer.json` → `require` and `require-dev`
- Skips `php` and `ext-*` pseudo-packages

## Risk Scoring Logic

```
risk_score = 0  if package exists on registry
risk_score = 100 if package NOT found (404) ← TAKEOVER TARGET
```

### Signals

- `not_found_on_npm` - 404 on npm registry
- `not_found_on_pypi` - 404 on PyPI
- `not_found_on_packagist` - 404 on Packagist

## Report Format

```json
{
  "repo_path": "string",
  "github_repo": "owner/repo",
  "ecosystems": {
    "npm": { ... },
    "pypi": { ... },
    "composer": { ... }
  }
}
```

Each ecosystem contains:
- `dependencies_by_file` - Map of files to packages
- `total_dependencies` - Count
- `risk_analysis` - Per-package results
- `summary` - High-level statistics

## Next Steps / Future Enhancements

- [ ] Web UI dashboard
- [ ] Parallel registry checks (goroutines)
- [ ] Caching of registry responses
- [ ] Download metrics enrichment
- [ ] Scheduled scanning for organizations
- [ ] Integration with bug bounty platforms
- [ ] Custom registry support
- [ ] Typosquatting detection
- [ ] Dependency graph visualization

---

**Built for speed, designed for hunters. 🎯**
