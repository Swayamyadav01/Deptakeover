# DepTakeover

```
 ____           _____     _
|  _ \  ___ _ _|_   _|_ _| | _____  _____   _____ _ __ 
| | | |/ _ \ '_ \| |/ _` | |/ / _ \/ _ \ \ / / _ \ '__|
| |_| |  __/ |_) | | (_| |   <  __/ (_) \ V /  __/ |   
|____/ \___| .__/|_|\__,_|_|\_\___|\___/ \_/ \___|_|   
           |_|

Supply Chain Takeover Scanner
Find missing packages across npm, PyPI, and Composer
Report unclaimed dependencies before attackers do
```

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/Swayamyadav01/Deptakeover/releases)

Package takeover finder for bug bounty hunting. Scans npm, PyPI, and Composer registries for unclaimed packages.

Built this after manually checking dependencies got old real fast during bug bounty hunts. Figured other people might find it useful too.

## What's package takeover?

Basically when a project depends on a package that doesn't exist anymore on the registry. You could potentially claim that package name and own everyone who depends on it. Pretty bad for supply chain security.

This tool finds those missing packages automatically instead of you having to check each one by hand.

## Features

- Scans npm, PyPI, and Composer registries
- Can scan entire GitHub organizations (this is where it gets useful)
- Pretty fast - written in Go with concurrent requests
- Outputs JSON reports for further analysis
- Works on Windows/Linux/macOS

## How it works

Simple - grabs dependencies from package files (package.json, requirements.txt, composer.json) then hits the registry APIs to check if they return 404. Those 404s are your potential takeover targets.

## Installation

```bash
# Go install - simplest way
go install github.com/Swayamyadav01/Deptakeover/cmd/deptakeover@latest
```

Or grab a binary from releases if you don't have Go installed.

```bash
# Build yourself
git clone https://github.com/Swayamyadav01/Deptakeover.git
cd Deptakeover  
go build -o deptakeover ./cmd/deptakeover
```

## Usage

Scan a single repo:
```bash
deptakeover npm facebook/react
deptakeover pypi django/django  
deptakeover composer laravel/laravel

# shortcuts
deptakeover py some/repo    # same as pypi
deptakeover php vendor/pkg  # same as composer
```

Scan entire organizations (this is where it gets interesting):
```bash
deptakeover org microsoft           # all repos, all package types
deptakeover org-npm facebook        # just npm
deptakeover org-pypi google         # just python
deptakeover org-composer symfony    # just php
```

## Example

```bash
$ deptakeover npm some/repo

Scanning npm dependencies...
Checking 47 packages...
Found 3 unclaimed packages!

Results saved to: npm_report.json
```

The JSON report has details about which packages returned 404.

## Project structure

```
cmd/deptakeover/      - main CLI app
internal/scanner/     - parses package.json, requirements.txt, etc
internal/registry/    - checks npm/pypi/packagist APIs
internal/github/      - GitHub repo cloning/downloading
scripts/              - build and release scripts
.github/              - GitHub workflows
build/                - build outputs (generated)
```

Pretty standard Go layout. The scanner modules find dependencies, registry modules check if they exist.

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

## Development

```bash
git clone https://github.com/Swayamyadav01/Deptakeover.git
cd Deptakeover  
go build ./cmd/deptakeover
```

Run tests:
```bash
go test ./...
```

## License

MIT - do whatever you want with it.

## Security

If you find bugs in this tool, just open an issue.

## Bug bounty tips

- Start with org-npm scans on JS-heavy companies - usually more dependencies
- Check if missing package names are typos of popular packages - those pay well
- Some packages get claimed/unclaimed over time, so re-scan targets periodically  
- Don't actually register packages - just report the potential takeover
- Always follow responsible disclosure

## TODO

- Add more registries (RubyGems, NuGet maybe)
- Better rate limit handling  
- Cache results to avoid re-scanning same repos
- Web interface if anyone wants that

## Contributing

Pull requests welcome. Keep it simple.

To add a new registry, look at existing ones in `internal/registry/` and follow the same pattern.

---


Built for fellow bug bounty hunters who got tired of manually checking dependencies. Hope it helps you find some good stuff.


