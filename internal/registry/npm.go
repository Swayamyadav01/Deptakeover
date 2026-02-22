package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NPMPackageInfo struct {
	Exists    bool
	RiskScore int
	Signals   []string
	Metadata  map[string]interface{}
	Package   string
}

func CheckNPMPackageRisk(packageName string) NPMPackageInfo {
	result := NPMPackageInfo{
		Package:   packageName,
		Exists:    false,
		RiskScore: 0,
		Signals:   []string{},
		Metadata:  make(map[string]interface{}),
	}

	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
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
		fmt.Printf("Warning: Error fetching %s from npm registry: %v\n", packageName, err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("Info: Package not found on npm: %s\n", packageName)
		result.RiskScore = 100
		result.Signals = []string{"not_found_on_npm"}
		return result
	}

	if resp.StatusCode == 200 {
		result.Exists = true
		result.RiskScore = 0

		body, _ := io.ReadAll(resp.Body)
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			result.Metadata = map[string]interface{}{
				"name":        data["name"],
				"description": data["description"],
				"repository":  "npm",
			}
		}
		return result
	}

	fmt.Printf("Warning: Unexpected status %d for %s\n", resp.StatusCode, packageName)
	return result
}

func AnalyzeNPMDependencyRisks(packages []string) map[string]NPMPackageInfo {
	results := make(map[string]NPMPackageInfo)
	for _, pkg := range packages {
		fmt.Printf("Analyzing %s...\n", pkg)
		results[pkg] = CheckNPMPackageRisk(pkg)
	}
	return results
}
