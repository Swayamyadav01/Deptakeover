package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FindPythonDependencyFiles(repoPath string) []string {
	var depFiles []string

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip venv, virtualenv, .venv, site-packages
		if info.IsDir() {
			if strings.HasSuffix(info.Name(), "venv") || strings.HasSuffix(info.Name(), "env") ||
				info.Name() == ".venv" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
		}

		// Check for Python dependency files
		if strings.HasSuffix(info.Name(), "requirements.txt") ||
			info.Name() == "setup.py" ||
			info.Name() == "pyproject.toml" ||
			info.Name() == "Pipfile" {
			depFiles = append(depFiles, path)
			fmt.Printf("Found Python dependency file: %s\n", path)
		}

		return nil
	})

	return depFiles
}

func ParseRequirementsText(filePath string) []string {
	var packages []string

	file, err := os.Open(filePath)
	if err != nil {
		return packages
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^([a-zA-Z0-9\-_.]+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Skip URL entries
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			continue
		}

		// Extract package name (before ==, >=, <=, etc.)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			pkgName := strings.ToLower(matches[1])
			if pkgName != "-e" && pkgName != "." {
				packages = append(packages, pkgName)
			}
		}
	}

	return packages
}

func ParseSetupPy(filePath string) []string {
	var packages []string

	data, err := os.ReadFile(filePath)
	if err != nil {
		return packages
	}

	contentStr := string(data)

	// Match install_requires = [...]  or requires = [...]
	re := regexp.MustCompile(`(?:install_requires|requires)\s*=\s*\[(.*?)\]`)
	matches := re.FindStringSubmatch(contentStr)

	if len(matches) > 1 {
		content := matches[1]
		// Extract package names from strings
		pkgRe := regexp.MustCompile(`["']([a-zA-Z0-9\-_]+)`)
		for _, match := range pkgRe.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				packages = append(packages, strings.ToLower(match[1]))
			}
		}
	}

	return packages
}

func ParsePyprojectToml(filePath string) []string {
	var packages []string

	data, err := os.ReadFile(filePath)
	if err != nil {
		return packages
	}

	contentStr := string(data)

	// Match dependencies section
	re := regexp.MustCompile(`dependencies\s*=\s*\[(.*?)\]`)
	matches := re.FindStringSubmatch(contentStr)

	if len(matches) > 1 {
		content := matches[1]
		// Extract package names
		pkgRe := regexp.MustCompile(`["']([a-zA-Z0-9\-_]+)`)
		for _, match := range pkgRe.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				packages = append(packages, strings.ToLower(match[1]))
			}
		}
	}

	return packages
}

func ParsePipfile(filePath string) []string {
	var packages []string

	data, err := os.ReadFile(filePath)
	if err != nil {
		return packages
	}

	contentStr := string(data)

	// Match packages section
	re := regexp.MustCompile(`\[packages\](.*?)(?:\[|$)`)
	matches := re.FindStringSubmatch(contentStr)

	if len(matches) > 1 {
		content := matches[1]
		// Extract package names
		pkgRe := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)\s*=`)
		for _, line := range strings.Split(content, "\n") {
			match := pkgRe.FindStringSubmatch(line)
			if len(match) > 1 {
				packages = append(packages, strings.ToLower(match[1]))
			}
		}
	}

	return packages
}

func ExtractAllPythonDependencies(repoPath string) map[string][]string {
	depsByFile := make(map[string][]string)

	depFiles := FindPythonDependencyFiles(repoPath)

	for _, depFile := range depFiles {
		var deps []string

		switch {
		case strings.HasSuffix(depFile, "requirements.txt"):
			deps = ParseRequirementsText(depFile)
		case strings.HasSuffix(depFile, "setup.py"):
			deps = ParseSetupPy(depFile)
		case strings.HasSuffix(depFile, "pyproject.toml"):
			deps = ParsePyprojectToml(depFile)
		case strings.HasSuffix(depFile, "Pipfile"):
			deps = ParsePipfile(depFile)
		}

		if len(deps) > 0 {
			depsByFile[depFile] = deps
		}
	}

	return depsByFile
}

func GetAllUniquePythonDeps(repoPath string) []string {
	depsByFile := ExtractAllPythonDependencies(repoPath)

	uniqueMap := make(map[string]bool)
	for _, deps := range depsByFile {
		for _, name := range deps {
			uniqueMap[name] = true
		}
	}

	var uniqueDeps []string
	for name := range uniqueMap {
		uniqueDeps = append(uniqueDeps, name)
	}

	return uniqueDeps
}
