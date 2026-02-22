package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Swayamyadav01/Deptakeover/internal/github"
	"github.com/Swayamyadav01/Deptakeover/internal/registry"
	"github.com/Swayamyadav01/Deptakeover/internal/scanner"

	"github.com/spf13/cobra"
)

type ReportData struct {
	RepoPath   string                   `json:"repo_path"`
	GitHubOrg  *string                  `json:"github_org"`
	GitHubRepo *string                  `json:"github_repo"`
	GitHubURL  *string                  `json:"github_url"`
	Ecosystems map[string]EcosystemData `json:"ecosystems"`
}

type EcosystemData struct {
	DependenciesByFile map[string]interface{} `json:"dependencies_by_file"`
	TotalDependencies  int                    `json:"total_dependencies"`
	RiskAnalysis       map[string]interface{} `json:"risk_analysis"`
	Summary            SummaryData            `json:"summary"`
}

type SummaryData struct {
	HighRiskCount    int      `json:"high_risk_count"`
	MediumRiskCount  int      `json:"medium_risk_count"`
	NotFoundCount    int      `json:"not_found_count"`
	HighRiskPackages []string `json:"high_risk_packages"`
	NotFoundPackages []string `json:"not_found_packages"`
}

var rootCmd = &cobra.Command{
	Use:   "deptakeover <ecosystem> <repo|org>",
	Short: "Package takeover scanner for bug bounty hunting",
	Long: `DepTakeover - Supply Chain Vulnerability Scanner

Hunt for unclaimed packages across npm, PyPI, and Composer registries.
When a project depends on a package that no longer exists on the registry,
an attacker can claim that package name and potentially compromise all 
projects that depend on it.

DepTakeover scans repositories and organizations to find these dangerous
missing dependencies automatically.

SINGLE REPOSITORY SCANNING:
  deptakeover npm lodash/lodash               # Scan npm dependencies
  deptakeover pypi django/django              # Scan Python packages
  deptakeover composer laravel/laravel        # Scan PHP packages

ORGANIZATION-WIDE SCANNING:
  deptakeover org microsoft                   # All ecosystems
  deptakeover org-npm facebook                # npm only
  deptakeover org-pypi google                 # PyPI only
  deptakeover org-composer symfony            # Composer only

SHORTCUTS:
  py = pypi, php = composer

OUTPUT:
  JSON report with all vulnerable dependencies and risk analysis

PERFECT FOR:
  Bug bounty hunters, security researchers, and DevOps teams looking
  to identify supply chain attack vectors in their dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(bannerText)
		fmt.Println()

		if len(args) != 2 {
			fmt.Println("Usage: deptakeover <ecosystem> <target>")
			fmt.Println()
			fmt.Println("SINGLE REPOSITORY:")
			fmt.Println("  deptakeover npm lodash/lodash")
			fmt.Println("  deptakeover pypi django/django")
			fmt.Println("  deptakeover composer laravel/laravel")
			fmt.Println("  deptakeover py requests               # shorthand")
			fmt.Println()
			fmt.Println("ORGANIZATION SCANNING:")
			fmt.Println("  deptakeover org microsoft              # All ecosystems")
			fmt.Println("  deptakeover org-npm facebook           # npm only")
			fmt.Println("  deptakeover org-pypi google            # PyPI only")
			fmt.Println("  deptakeover org-composer symfony       # Composer only")
			fmt.Println()
			fmt.Println("OUTPUT:")
			fmt.Println("  Generates JSON report with vulnerable packages")
			fmt.Println()
			fmt.Println("Need help? Run: deptakeover --help")
			os.Exit(1)
		}

		ecosystemInput := args[0]
		targetInput := args[1]

		// Map ecosystem names
		ecosystemMap := map[string]string{
			"npm":          "npm",
			"pypi":         "pypi",
			"py":           "pypi",
			"python":       "pypi",
			"composer":     "composer",
			"php":          "composer",
			"org":          "org",
			"org-npm":      "org-npm",
			"org-pypi":     "org-pypi",
			"org-composer": "org-composer",
		}

		ecosystem, exists := ecosystemMap[ecosystemInput]
		if !exists {
			fmt.Printf("Unknown ecosystem: '%s'\n", ecosystemInput)
			fmt.Println("Valid options: npm, pypi, py, composer, php")
			fmt.Println("Org scans: org, org-npm, org-pypi, org-composer")
			fmt.Println("Example: deptakeover npm lodash/lodash")
			os.Exit(1)
		}

		// Handle organization scanning
		if strings.HasPrefix(ecosystem, "org") {
			runOrgScan(ecosystem, targetInput)
			return
		}

		// Handle single repository scanning
		outFile := ecosystemInput + "_report.json"

		// Determine if GitHub URL or owner/repo
		var githubRepo, githubURL string
		if strings.HasPrefix(targetInput, "https://") || strings.HasPrefix(targetInput, "http://") || strings.HasPrefix(targetInput, "git@") {
			githubURL = targetInput
		} else {
			githubRepo = targetInput
		}

		runScan(githubURL, githubRepo, "", "", ecosystem, outFile)
	},
}

const bannerText = `
  ____           _____     _
|  _ \  ___ _ _|_   _|_ _| | _____  _____   _____ _ __ 
| | | |/ _ \ '_ \| |/ _' | |/ / _ \/ _ \ \ / / _ \ '__|
| |_| |  __/ |_) | | (_| |   <  __/ (_) \ V /  __/ |   
|____/ \___| .__/|_|\__,_|_|\_\___|\___/ \_/ \___|_|   
           |_| 
		   
DepTakeover
Supply Chain Takeover Scanner
Find missing packages across npm, PyPI, and Composer
Report unclaimed dependencies before attackers do
`

func init() {
	// Positional arguments only - no flags
}

func runScan(githubURL, githubRepo, githubOrg, localPath, ecosystem, outFile string) {
	fmt.Printf("🔍 Scanning [%s]...\n", ecosystem)

	// Get repo path
	repoPath, err := github.GetRepoPath(githubURL, githubRepo, githubOrg, localPath)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	report := ReportData{
		RepoPath:   repoPath,
		Ecosystems: make(map[string]EcosystemData),
	}

	if githubRepo != "" {
		report.GitHubRepo = &githubRepo
	}
	if githubURL != "" {
		report.GitHubURL = &githubURL
	}

	// Single ecosystem scan
	if ecosystem == "npm" {
		depsByFile := scanner.ExtractAllNPMDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniqueNPMDeps(repoPath)
			fmt.Printf("📦 Found %d packages\n", len(allDeps))

			riskAnalysis := registry.AnalyzeNPMDependencyRisks(allDeps)

			npm := EcosystemData{
				DependenciesByFile: convertNPMDeps(depsByFile),
				TotalDependencies:  len(allDeps),
				RiskAnalysis:       convertNPMAnalysis(riskAnalysis),
				Summary:            generateSummary(riskAnalysis),
			}
			report.Ecosystems["npm"] = npm
		}
	}

	if ecosystem == "pypi" {
		depsByFile := scanner.ExtractAllPythonDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniquePythonDeps(repoPath)
			fmt.Printf("📦 Found %d packages\n", len(allDeps))

			riskAnalysis := registry.AnalyzePyPIDependencyRisks(allDeps)

			pypi := EcosystemData{
				DependenciesByFile: convertPythonDeps(depsByFile),
				TotalDependencies:  len(allDeps),
				RiskAnalysis:       convertPyPIAnalysis(riskAnalysis),
				Summary:            generatePyPISummary(riskAnalysis),
			}
			report.Ecosystems["pypi"] = pypi
		}
	}

	if ecosystem == "composer" {
		depsByFile := scanner.ExtractAllPHPDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniquePHPDeps(repoPath)
			fmt.Printf("📦 Found %d packages\n", len(allDeps))

			riskAnalysis := registry.AnalyzePackagistDependencyRisks(allDeps)

			composer := EcosystemData{
				DependenciesByFile: convertPHPDeps(depsByFile),
				TotalDependencies:  len(allDeps),
				RiskAnalysis:       convertPackagistAnalysis(riskAnalysis),
				Summary:            generatePackagistSummary(riskAnalysis),
			}
			report.Ecosystems["composer"] = composer
		}
	}

	if len(report.Ecosystems) == 0 {
		fmt.Println("⚠️  No dependencies found")
		return
	}

	// Save report
	os.MkdirAll(filepath.Dir(outFile), 0755)
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(outFile, reportJSON, 0644)

	fmt.Printf("✅ Report: %s\n", outFile)
	printSummary(report)
}

// GitHub API response structure for repository listing
type GitHubRepo struct {
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	Language   string `json:"language"`
	Size       int    `json:"size"`
	ForksCount int    `json:"forks_count"`
	StarsCount int    `json:"stargazers_count"`
}

// Organization scan report structure
type OrgReportData struct {
	Organization       string                    `json:"organization"`
	ScanType           string                    `json:"scan_type"`
	TotalRepos         int                       `json:"total_repos"`
	ScannedRepos       int                       `json:"scanned_repos"`
	SkippedRepos       int                       `json:"skipped_repos"`
	TotalVulns         int                       `json:"total_vulnerabilities"`
	RepositorySummary  map[string]RepoScanResult `json:"repository_summary"`
	TopVulnerabilities []VulnSummary             `json:"top_vulnerabilities"`
	ScanTimestamp      time.Time                 `json:"scan_timestamp"`
}

type RepoScanResult struct {
	Language     string   `json:"language,omitempty"`
	Stars        int      `json:"stars"`
	Size         int      `json:"size_kb"`
	VulnCount    int      `json:"vulnerability_count"`
	VulnPackages []string `json:"vulnerable_packages"`
	ScanStatus   string   `json:"scan_status"`
	Error        string   `json:"error,omitempty"`
}

type VulnSummary struct {
	PackageName  string   `json:"package_name"`
	Ecosystem    string   `json:"ecosystem"`
	FoundInRepos []string `json:"found_in_repos"`
	Frequency    int      `json:"frequency"`
}

func runOrgScan(scanType, orgName string) {
	fmt.Printf("🔍 Organization Scan: %s [%s]\n", orgName, scanType)

	// Get repositories from GitHub API
	repos, err := getOrgRepositories(orgName)
	if err != nil {
		fmt.Printf("❌ Error fetching repositories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📁 Found %d repositories\n", len(repos))

	report := OrgReportData{
		Organization:      orgName,
		ScanType:          scanType,
		TotalRepos:        len(repos),
		RepositorySummary: make(map[string]RepoScanResult),
		ScanTimestamp:     time.Now(),
	}

	vulnMap := make(map[string]*VulnSummary)

	// Determine which ecosystems to scan
	ecosystems := getEcosystemsForOrgScan(scanType)

	count := 0
	for _, repo := range repos {
		count++
		fmt.Printf("\n[%d/%d] Scanning %s...\n", count, len(repos), repo.FullName)

		// Skip very large repos (> 100MB) to avoid timeouts
		if repo.Size > 100000 {
			fmt.Printf("⏭️  Skipped (too large: %d KB)\n", repo.Size)
			report.SkippedRepos++
			report.RepositorySummary[repo.Name] = RepoScanResult{
				Language:   repo.Language,
				Stars:      repo.StarsCount,
				Size:       repo.Size,
				ScanStatus: "skipped_large",
			}
			continue
		}

		repoResult := RepoScanResult{
			Language:   repo.Language,
			Stars:      repo.StarsCount,
			Size:       repo.Size,
			ScanStatus: "scanned",
		}

		// Run scans for each ecosystem
		for _, eco := range ecosystems {
			vulnPackages, err := scanRepoForEcosystem(repo.FullName, eco)
			if err != nil {
				repoResult.Error = err.Error()
				repoResult.ScanStatus = "error"
				continue
			}

			// Collect vulnerabilities
			for _, pkg := range vulnPackages {
				repoResult.VulnPackages = append(repoResult.VulnPackages, fmt.Sprintf("%s:%s", eco, pkg))
				repoResult.VulnCount++

				// Track for overall summary
				key := eco + ":" + pkg
				if vuln, exists := vulnMap[key]; exists {
					vuln.Frequency++
					vuln.FoundInRepos = append(vuln.FoundInRepos, repo.Name)
				} else {
					vulnMap[key] = &VulnSummary{
						PackageName:  pkg,
						Ecosystem:    eco,
						FoundInRepos: []string{repo.Name},
						Frequency:    1,
					}
				}
			}
		}

		report.RepositorySummary[repo.Name] = repoResult
		report.ScannedRepos++
		report.TotalVulns += repoResult.VulnCount

		// Rate limiting - pause between repos
		time.Sleep(500 * time.Millisecond)
	}

	// Generate top vulnerabilities list
	for _, vuln := range vulnMap {
		report.TopVulnerabilities = append(report.TopVulnerabilities, *vuln)
	}

	// Sort by frequency (most common vulnerabilities first)
	sort.Slice(report.TopVulnerabilities, func(i, j int) bool {
		return report.TopVulnerabilities[i].Frequency > report.TopVulnerabilities[j].Frequency
	})

	// Save organization report
	outFile := fmt.Sprintf("%s_%s_report.json", orgName, scanType)
	os.MkdirAll(filepath.Dir(outFile), 0755)
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(outFile, reportJSON, 0644)

	// Print summary
	printOrgSummary(report)
}

func getOrgRepositories(orgName string) ([]GitHubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?type=public&per_page=100", orgName)
	var allRepos []GitHubRepo
	page := 1

	for {
		pageURL := fmt.Sprintf("%s&page=%d", url, page)
		resp, err := http.Get(pageURL)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("organization '%s' not found", orgName)
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("API error: %d", resp.StatusCode)
		}

		var repos []GitHubRepo
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, fmt.Errorf("JSON decode error: %w", err)
		}

		if len(repos) == 0 {
			break // No more pages
		}

		allRepos = append(allRepos, repos...)
		page++

		// Respect rate limits
		time.Sleep(100 * time.Millisecond)
	}

	return allRepos, nil
}

func getEcosystemsForOrgScan(scanType string) []string {
	switch scanType {
	case "org":
		return []string{"npm", "pypi", "composer"}
	case "org-npm":
		return []string{"npm"}
	case "org-pypi":
		return []string{"pypi"}
	case "org-composer":
		return []string{"composer"}
	default:
		return []string{"npm", "pypi", "composer"}
	}
}

func scanRepoForEcosystem(repoFullName, ecosystem string) ([]string, error) {
	// Clone or use cached repo
	repoPath, err := github.GetRepoPath("", repoFullName, "", "")
	if err != nil {
		return nil, err
	}

	var vulnerablePackages []string

	switch ecosystem {
	case "npm":
		depsByFile := scanner.ExtractAllNPMDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniqueNPMDeps(repoPath)
			riskAnalysis := registry.AnalyzeNPMDependencyRisks(allDeps)
			for pkg, risk := range riskAnalysis {
				if risk.RiskScore >= 70 || !risk.Exists {
					vulnerablePackages = append(vulnerablePackages, pkg)
				}
			}
		}
	case "pypi":
		depsByFile := scanner.ExtractAllPythonDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniquePythonDeps(repoPath)
			riskAnalysis := registry.AnalyzePyPIDependencyRisks(allDeps)
			for pkg, risk := range riskAnalysis {
				if risk.RiskScore >= 70 || !risk.Exists {
					vulnerablePackages = append(vulnerablePackages, pkg)
				}
			}
		}
	case "composer":
		depsByFile := scanner.ExtractAllPHPDependencies(repoPath)
		if len(depsByFile) > 0 {
			allDeps := scanner.GetAllUniquePHPDeps(repoPath)
			riskAnalysis := registry.AnalyzePackagistDependencyRisks(allDeps)
			for pkg, risk := range riskAnalysis {
				if risk.RiskScore >= 70 || !risk.Exists {
					vulnerablePackages = append(vulnerablePackages, pkg)
				}
			}
		}
	}

	return vulnerablePackages, nil
}

func printOrgSummary(report OrgReportData) {
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("🏢 ORGANIZATION SCAN SUMMARY: %s\n", strings.ToUpper(report.Organization))
	fmt.Println(strings.Repeat("═", 60))

	fmt.Printf("📊 Scan Type: %s\n", report.ScanType)
	fmt.Printf("📁 Repositories: %d total, %d scanned, %d skipped\n",
		report.TotalRepos, report.ScannedRepos, report.SkippedRepos)
	fmt.Printf("🚨 Total Vulnerabilities: %d\n", report.TotalVulns)

	if len(report.TopVulnerabilities) > 0 {
		fmt.Printf("\n🔥 TOP VULNERABLE PACKAGES:\n")
		count := 0
		for _, vuln := range report.TopVulnerabilities {
			if count >= 10 { // Show top 10
				break
			}
			fmt.Printf("  %d. %s [%s] - found in %d repos\n",
				count+1, vuln.PackageName, vuln.Ecosystem, vuln.Frequency)
			count++
		}
	}

	vulnerable := 0
	for _, repo := range report.RepositorySummary {
		if repo.VulnCount > 0 {
			vulnerable++
		}
	}

	fmt.Printf("\n📈 REPOSITORIES WITH VULNERABILITIES: %d/%d\n", vulnerable, report.ScannedRepos)

	// Show most vulnerable repos
	type repoVulnCount struct {
		name  string
		count int
	}
	var repoVulns []repoVulnCount
	for name, repo := range report.RepositorySummary {
		if repo.VulnCount > 0 {
			repoVulns = append(repoVulns, repoVulnCount{name, repo.VulnCount})
		}
	}

	sort.Slice(repoVulns, func(i, j int) bool {
		return repoVulns[i].count > repoVulns[j].count
	})

	if len(repoVulns) > 0 {
		fmt.Println("\n🎯 MOST VULNERABLE REPOSITORIES:")
		count := 0
		for _, rv := range repoVulns {
			if count >= 5 { // Show top 5
				break
			}
			fmt.Printf("  %d. %s (%d vulnerabilities)\n", count+1, rv.name, rv.count)
			count++
		}
	}

	fmt.Printf("\n✅ Report saved: %s_%s_report.json\n", report.Organization, report.ScanType)
	fmt.Println(strings.Repeat("═", 60))
}

func generateSummary(riskAnalysis map[string]registry.NPMPackageInfo) SummaryData {
	var highRisk, mediumRisk, notFound []string

	for pkg, risk := range riskAnalysis {
		if risk.RiskScore >= 70 {
			highRisk = append(highRisk, pkg)
		} else if risk.RiskScore >= 40 {
			mediumRisk = append(mediumRisk, pkg)
		} else if !risk.Exists {
			notFound = append(notFound, pkg)
		}
	}

	sort.Strings(highRisk)
	sort.Strings(mediumRisk)
	sort.Strings(notFound)

	if len(highRisk) > 20 {
		highRisk = highRisk[:20]
	}
	if len(notFound) > 20 {
		notFound = notFound[:20]
	}

	return SummaryData{
		HighRiskCount:    len(highRisk),
		MediumRiskCount:  len(mediumRisk),
		NotFoundCount:    len(notFound),
		HighRiskPackages: highRisk,
		NotFoundPackages: notFound,
	}
}

func generatePyPISummary(riskAnalysis map[string]registry.PyPIPackageInfo) SummaryData {
	var highRisk, mediumRisk, notFound []string

	for pkg, risk := range riskAnalysis {
		if risk.RiskScore >= 70 {
			highRisk = append(highRisk, pkg)
		} else if risk.RiskScore >= 40 {
			mediumRisk = append(mediumRisk, pkg)
		} else if !risk.Exists {
			notFound = append(notFound, pkg)
		}
	}

	sort.Strings(highRisk)
	sort.Strings(mediumRisk)
	sort.Strings(notFound)

	if len(highRisk) > 20 {
		highRisk = highRisk[:20]
	}
	if len(notFound) > 20 {
		notFound = notFound[:20]
	}

	return SummaryData{
		HighRiskCount:    len(highRisk),
		MediumRiskCount:  len(mediumRisk),
		NotFoundCount:    len(notFound),
		HighRiskPackages: highRisk,
		NotFoundPackages: notFound,
	}
}

func generatePackagistSummary(riskAnalysis map[string]registry.PackagistPackageInfo) SummaryData {
	var highRisk, mediumRisk, notFound []string

	for pkg, risk := range riskAnalysis {
		if risk.RiskScore >= 70 {
			highRisk = append(highRisk, pkg)
		} else if risk.RiskScore >= 40 {
			mediumRisk = append(mediumRisk, pkg)
		} else if !risk.Exists {
			notFound = append(notFound, pkg)
		}
	}

	sort.Strings(highRisk)
	sort.Strings(mediumRisk)
	sort.Strings(notFound)

	if len(highRisk) > 20 {
		highRisk = highRisk[:20]
	}
	if len(notFound) > 20 {
		notFound = notFound[:20]
	}

	return SummaryData{
		HighRiskCount:    len(highRisk),
		MediumRiskCount:  len(mediumRisk),
		NotFoundCount:    len(notFound),
		HighRiskPackages: highRisk,
		NotFoundPackages: notFound,
	}
}

func convertNPMAnalysis(analysis map[string]registry.NPMPackageInfo) map[string]interface{} {
	result := make(map[string]interface{})
	for pkg, risk := range analysis {
		result[pkg] = risk
	}
	return result
}

func convertNPMDeps(deps map[string]map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for file, dep := range deps {
		result[file] = dep
	}
	return result
}

func convertPyPIAnalysis(analysis map[string]registry.PyPIPackageInfo) map[string]interface{} {
	result := make(map[string]interface{})
	for pkg, risk := range analysis {
		result[pkg] = risk
	}
	return result
}

func convertPackagistAnalysis(analysis map[string]registry.PackagistPackageInfo) map[string]interface{} {
	result := make(map[string]interface{})
	for pkg, risk := range analysis {
		result[pkg] = risk
	}
	return result
}

func convertPythonDeps(deps map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for file, pkgs := range deps {
		result[file] = pkgs
	}
	return result
}

func convertPHPDeps(deps map[string]map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for file, dep := range deps {
		result[file] = dep
	}
	return result
}

func printSummary(report ReportData) {
	fmt.Println(strings.Repeat("─", 50))

	totalDeps := 0
	totalNotFound := 0

	for _, eco := range report.Ecosystems {
		totalDeps += eco.TotalDependencies
		totalNotFound += eco.Summary.NotFoundCount
	}

	fmt.Printf("📊 Dependencies: %d\n", totalDeps)
	fmt.Printf("⚠️  Takeover targets: %d\n", totalNotFound)

	for ecoName, ecoData := range report.Ecosystems {
		if ecoData.Summary.NotFoundCount > 0 {
			fmt.Printf("\n🚨 [%s] %d NOT FOUND:\n", strings.ToUpper(ecoName), ecoData.Summary.NotFoundCount)
			for _, pkg := range ecoData.Summary.NotFoundPackages {
				fmt.Printf("  • %s\n", pkg)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
