# Quick Start - deptakeover

## Updated CLI Syntax ✨

The tool now uses **minimal positional arguments** instead of verbose flags:

### New Usage
```bash
deptakeover <ecosystem> <repo>
```

### Examples

**npm scanning:**
```bash
deptakeover npm lodash/lodash
deptakeover npm https://github.com/lodash/lodash
```

**PyPI scanning (Python):**
```bash
deptakeover pypi django/django
deptakeover py django/django                    # Shorthand
```

**Composer scanning (PHP):**
```bash
deptakeover composer laravel/laravel
deptakeover php laravel/laravel                 # Shorthand
```

### Installation

1. **Download Go** (if not already installed):
   - Visit: https://golang.org/dl/
   - Download Go 1.21 or later for your OS
   - Install and add to PATH

2. **Build deptakeover**:
   ```bash
   cd "c:\Bug Bounty\supply chain attack tool"
   .\build.bat                                   # Windows
   # OR
   chmod +x build.sh && ./build.sh              # Linux/macOS
   ```

3. **Run**:
   ```bash
   # Windows
   deptakeover.exe npm lodash/lodash
   npm_report.json will be generated
   
   # Linux/macOS
   ./deptakeover npm lodash/lodash
   npm_report.json will be generated
   ```

## What Changed?

### Old Syntax (removed)
```bash
deptakeover --github-repo lodash/lodash --ecosystem npm --out report.json
```

### New Syntax (current)
```bash
deptakeover npm lodash/lodash
```

### Features
- ✅ Positional arguments only (no flags)
- ✅ Ecosystem shorthand: `py` → PyPI, `php` → Composer  
- ✅ Auto-detect GitHub URL vs owner/repo format
- ✅ Output: `{ecosystem}_report.json` (npm_report.json, pypi_report.json, etc.)
- ✅ Minimal emoji-based console output
- ✅ Fast execution with parallel registry checks

## Setup Requirements

- **Go 1.21+** - For building
- **git** (optional) - For cloning GitHub repos; zip fallback available
- **No Python/PHP/Node.js required** - Standalone compiled binary

## Testing After Build

```bash
# Test on lodash (npm)
deptakeover npm lodash/lodash

# Expected output:
# 🔍 Scanning [npm]...
# 📦 Found 27 packages
# ✅ Report: npm_report.json
# 📊 Dependencies: 27
# ⚠️  Takeover targets: 0
```

## Troubleshooting

**"command not found: deptakeover"**
- Add `./build/` to PATH, or run executables with full path
- Windows: `.\build\deptakeover-windows-amd64.exe`
- Linux: `./build/deptakeover-linux-amd64`

**"git clone failed"**
- Tool automatically falls back to zip download from GitHub
- No git installation required on target system

**build.bat not found**
- Use: `.\build.bat` (not `build.bat`)

**"go: command not found"**
- Go is not installed. Download and install from https://golang.org/dl/
