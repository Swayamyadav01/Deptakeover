package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PyPIPackageInfo struct {
	Exists    bool
	RiskScore int
	Signals   []string
	Metadata  map[string]interface{}
	Package   string
}

func CheckPyPIPackageRisk(packageName string) PyPIPackageInfo {
	result := PyPIPackageInfo{
		Package:   packageName,
		Exists:    false,
		RiskScore: 0,
		Signals:   []string{},
		Metadata:  make(map[string]interface{}),
	}

	// Skip special entries like "-e ."
	if packageName == "-e ." {
		result.RiskScore = 100
		result.Signals = []string{"not_found_on_pypi"}
		return result
	}

	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)
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
		fmt.Printf("Warning: Error fetching %s from PyPI: %v\n", packageName, err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("Info: Package not found on PyPI: %s\n", packageName)
		result.RiskScore = 100
		result.Signals = []string{"not_found_on_pypi"}
		return result
	}

	if resp.StatusCode == 200 {
		result.Exists = true
		result.RiskScore = 0

		body, _ := io.ReadAll(resp.Body)
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			if info, ok := data["info"].(map[string]interface{}); ok {
				result.Metadata = map[string]interface{}{
					"name":        info["name"],
					"description": info["summary"],
					"repository":  info["home_page"],
				}
			}
		}
		return result
	}

	fmt.Printf("Warning: Unexpected status %d for %s\n", resp.StatusCode, packageName)
	return result
}

func AnalyzePyPIDependencyRisks(packages []string) map[string]PyPIPackageInfo {
	results := make(map[string]PyPIPackageInfo)
	for _, pkg := range packages {
		fmt.Printf("Analyzing %s...\n", pkg)
		results[pkg] = CheckPyPIPackageRisk(pkg)
	}
	return results
}
