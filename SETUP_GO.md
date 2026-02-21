# Quick Start Guide for deptakeover (Go Version)

## Step 1: Install Go

deptakeover is written in Go. You need Go 1.21 or later installed.

### Windows

1. Download Go from https://golang.org/dl/
2. Choose `Windows` → `Installer (.msi)` → Run installer
3. Use default installation path: `C:\Program Files\Go`
4. **Restart your terminal/PowerShell** after installation

Verify installation:
```powershell
go version
# Should show: go version go1.21.x windows/amd64
```

### macOS

```bash
brew install go
# or download from https://golang.org/dl/

go version
```

### Linux

```bash
# Ubuntu/Debian
sudo apt-get install golang-go

# macOS alternative
brew install go

# Verify
go version
```

## Step 2: Clone the Repository

```powershell
# Windows PowerShell
git clone https://github.com/yourusername/deptakeover.git
cd deptakeover
```

Or if you already have the files, navigate to the deptakeover directory:
```powershell
cd "c:\Bug Bounty\supply chain attack tool"
```

## Step 3: Download Dependencies

```bash
go mod download
```

## Step 4: Build the Binary

### Option A: Windows Batch Script (Recommended)

```powershell
# Windows PowerShell
.\build.bat
```

This builds for Windows, Linux, and macOS.

### Option B: Manual Build

```powershell
# Build for Windows only
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o deptakeover.exe ./cmd/deptakeover
```

### Option C: Using Makefile (Windows with Git Bash or WSL)

```bash
make build
```

## Step 5: Run deptakeover

### Windows
```powershell
# Single file
.\deptakeover.exe --github-repo lodash/lodash --out report.json

# Or from build directory
.\build\deptakeover-windows-amd64.exe --help
```

### Linux/macOS
```bash
./deptakeover-linux-amd64 --github-repo lodash/lodash --out report.json
# or
./deptakeover-macos-amd64 --github-repo lodash/lodash --out report.json
```

## Step 6: Add to PATH (Optional)

### Windows

1. Find the binary location (e.g., `C:\deptakeover\build\`)
2. Right-click **Start Menu** → **System** → **Advanced system settings**
3. Click **Environment Variables**
4. Under **User variables**, click **New**
   - Variable name: `PATH`
   - Variable value: `C:\deptakeover\build\`
5. Click **OK** and restart PowerShell

Now you can run anywhere:
```powershell
deptakeover.exe --help
```

### Linux/macOS

```bash
sudo mv deptakeover-linux-amd64 /usr/local/bin/deptakeover
sudo chmod +x /usr/local/bin/deptakeover

# Now run from anywhere
deptakeover --help
```

## Common Issues

### "go: command not found"
**Solution**: Go is not installed or not in PATH. Restart terminal after installing Go.

### "can't load package"
**Solution**: Download dependencies with `go mod download`

### Build errors
**Solution**: 
1. Check Go version: `go version` (should be 1.21+)
2. Make sure you're in the deptakeover directory
3. Run `go mod tidy` to fix dependencies

## Next Steps

```bash
# See all available commands
deptakeover --help

# Scan a GitHub repo
deptakeover --github-repo lodash/lodash --out report.json

# Scan a local project
deptakeover --local-path ./my-project --out report.json

# View the report
cat report.json | jq .
```

## Performance Tips

- **First run**: Initial dependencies download takes ~30 seconds
- **Subsequent runs**: Much faster (cached)
- **Large repos**: May take 1-2 minutes for 100+ dependencies

## Troubleshooting

If you encounter issues:

1. **Update Go**: `go get -u all`
2. **Clear cache**: `go clean -cache`
3. **Reinstall deps**: `go mod tidy && go mod download`
4. **Check logs**: Build with verbose: `go build -v ./cmd/deptakeover`

## Get Support

- Check the README_GO.md for detailed documentation
- Review example commands in README_GO.md
- Check GitHub Issues for known problems

---

**You're all set! 🎯 Happy hunting!**
