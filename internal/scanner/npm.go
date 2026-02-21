package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func FindPackageJSONs(repoPath string) []string {
	var packageJSONs []string

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip node_modules and hidden directories
		if info.IsDir() && (info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".")) {
			return filepath.SkipDir
		}

		if info.Name() == "package.json" {
			packageJSONs = append(packageJSONs, path)
			fmt.Printf("Found package.json: %s\n", path)
		}

		return nil
	})

	return packageJSONs
}

func ExtractNPMDependencies(packageJSONPath string) (map[string]string, error) {
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	allDeps := make(map[string]string)

	// Add production dependencies
	for name, version := range pkg.Dependencies {
		// Skip scoped packages
		if !strings.HasPrefix(name, "@") {
			allDeps[name] = version
		}
	}

	// Add dev dependencies
	for name, version := range pkg.DevDependencies {
		// Skip scoped packages
		if !strings.HasPrefix(name, "@") {
			allDeps[name] = version
		}
	}

	fmt.Printf("Parsed %s: %d total dependencies\n", packageJSONPath, len(allDeps))
	return allDeps, nil
}

func ExtractAllNPMDependencies(repoPath string) map[string]map[string]string {
	depsByFile := make(map[string]map[string]string)

	packageJSONs := FindPackageJSONs(repoPath)

	for _, packageJSONPath := range packageJSONs {
		deps, err := ExtractNPMDependencies(packageJSONPath)
		if err == nil && len(deps) > 0 {
			depsByFile[packageJSONPath] = deps
		}
	}

	return depsByFile
}

func GetAllUniqueNPMDeps(repoPath string) []string {
	depsByFile := ExtractAllNPMDependencies(repoPath)

	uniqueMap := make(map[string]bool)
	for _, deps := range depsByFile {
		for name := range deps {
			uniqueMap[name] = true
		}
	}

	var uniqueDeps []string
	for name := range uniqueMap {
		uniqueDeps = append(uniqueDeps, name)
	}

	return uniqueDeps
}
