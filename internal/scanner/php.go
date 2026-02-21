package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ComposerJSON struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

func FindComposerJSONs(repoPath string) []string {
	var composerFiles []string

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip vendor and node_modules directories
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
		}

		if info.Name() == "composer.json" {
			composerFiles = append(composerFiles, path)
			fmt.Printf("Found composer.json: %s\n", path)
		}

		return nil
	})

	return composerFiles
}

func ExtractPHPDependencies(composerJSONPath string) map[string][]string {
	result := map[string][]string{
		"require":     {},
		"require-dev": {},
	}

	data, err := os.ReadFile(composerJSONPath)
	if err != nil {
		return result
	}

	var composer ComposerJSON
	if err := json.Unmarshal(data, &composer); err != nil {
		return result
	}

	// Process require dependencies
	for pkgName := range composer.Require {
		// Skip PHP core and extensions
		if !strings.HasPrefix(pkgName, "php") && !strings.HasPrefix(pkgName, "ext-") && strings.Contains(pkgName, "/") {
			result["require"] = append(result["require"], strings.ToLower(pkgName))
		}
	}

	// Process require-dev dependencies
	for pkgName := range composer.RequireDev {
		// Skip PHP core and extensions
		if !strings.HasPrefix(pkgName, "php") && !strings.HasPrefix(pkgName, "ext-") && strings.Contains(pkgName, "/") {
			result["require-dev"] = append(result["require-dev"], strings.ToLower(pkgName))
		}
	}

	fmt.Printf("Parsed %s: %d require, %d require-dev\n", composerJSONPath, len(result["require"]), len(result["require-dev"]))
	return result
}

func ExtractAllPHPDependencies(repoPath string) map[string]map[string][]string {
	depsByFile := make(map[string]map[string][]string)

	composerFiles := FindComposerJSONs(repoPath)

	for _, composerFile := range composerFiles {
		deps := ExtractPHPDependencies(composerFile)
		if len(deps["require"]) > 0 || len(deps["require-dev"]) > 0 {
			depsByFile[composerFile] = deps
		}
	}

	return depsByFile
}

func GetAllUniquePHPDeps(repoPath string) []string {
	depsByFile := ExtractAllPHPDependencies(repoPath)

	uniqueMap := make(map[string]bool)
	for _, deps := range depsByFile {
		for _, name := range deps["require"] {
			uniqueMap[name] = true
		}
		for _, name := range deps["require-dev"] {
			uniqueMap[name] = true
		}
	}

	var uniqueDeps []string
	for name := range uniqueMap {
		uniqueDeps = append(uniqueDeps, name)
	}

	return uniqueDeps
}
