@echo off
REM Build script for deptakeover on Windows

echo.
echo 🔨 Building deptakeover...
echo.

if not exist build mkdir build

echo 📦 Building for Windows (64-bit)...
set GOOS=windows
set GOARCH=amd64
go build -o build\deptakeover-windows-amd64.exe .\cmd\deptakeover

if %ERRORLEVEL% NEQ 0 (
    echo ❌ Build failed!
    exit /b 1
)

echo 📦 Building for Linux (64-bit)...
set GOOS=linux
set GOARCH=amd64
go build -o build\deptakeover-linux-amd64 .\cmd\deptakeover

echo 📦 Building for macOS (Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -o build\deptakeover-macos-amd64 .\cmd\deptakeover

echo 📦 Building for macOS (Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -o build\deptakeover-macos-arm64 .\cmd\deptakeover

set GOOS=
set GOARCH=

echo.
echo ✅ Build complete!
echo.
echo Binaries available in build\ folder:
echo   - deptakeover-windows-amd64.exe
echo   - deptakeover-linux-amd64
echo   - deptakeover-macos-amd64
echo   - deptakeover-macos-arm64
echo.
echo 🚀 Run with: .\build\deptakeover-windows-amd64.exe --help
echo.
