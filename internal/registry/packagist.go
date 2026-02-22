package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PackagistPackageInfo struct {
	Exists    bool
	RiskScore int
	Signals   []string
	Metadata  map[string]interface{}
	Package   string
}

type PackagistPackageJSON struct {
	Package struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Repository  string `json:"repository"`
	} `json:"package"`
}

func CheckPackagistPackageRisk(packageName string) PackagistPackageInfo {
	result := PackagistPackageInfo{
		Package:   packageName,
		Exists:    false,
		RiskScore: 0,
		Signals:   []string{},
		Metadata:  make(map[string]interface{}),
	}

	url := fmt.Sprintf("https://packagist.org/packages/%s.json", packageName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for %s: %v\n", packageName, err)
		return result
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: Error fetching %s from Packagist: %v\n", packageName, err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("Info: Package not found on Packagist: %s\n", packageName)
		result.RiskScore = 100
		result.Signals = []string{"not_found_on_packagist"}
		return result
	}

	if resp.StatusCode == 200 {
		result.Exists = true
		result.RiskScore = 0

		body, _ := io.ReadAll(resp.Body)
		var data PackagistPackageJSON
		if err := json.Unmarshal(body, &data); err == nil {
			result.Metadata = map[string]interface{}{
				"name":        data.Package.Name,
				"description": data.Package.Description,
				"repository":  data.Package.Repository,
			}
		}
		return result
	}

	fmt.Printf("Warning: Unexpected status %d for %s\n", resp.StatusCode, packageName)
	return result
}

func AnalyzePackagistDependencyRisks(packages []string) map[string]PackagistPackageInfo {
	results := make(map[string]PackagistPackageInfo)
	for _, pkg := range packages {
		fmt.Printf("Analyzing %s...\n", pkg)
		results[pkg] = CheckPackagistPackageRisk(pkg)
	}
	return results
}
