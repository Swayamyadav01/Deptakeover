package github

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func CloneGitHubRepo(repoURL, destPath string) error {
	fmt.Printf("Cloning %s to %s...\n", repoURL, destPath)

	// Try git clone first
	_, err := git.PlainClone(destPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err == nil {
		fmt.Printf("Successfully cloned via git.\n")
		return nil
	}

	fmt.Printf("Git clone failed: %v. Trying zip download...\n", err)

	// Fallback: Download as ZIP
	return downloadZip(repoURL, destPath)
}

func downloadZip(repoURL, destPath string) error {
	// Convert git URL to zip URL
	zipURL := strings.TrimSuffix(repoURL, ".git") + "/archive/refs/heads/main.zip"

	// Try main branch
	resp, err := http.Get(zipURL)
	if err != nil || resp.StatusCode == 404 {
		if resp != nil {
			resp.Body.Close()
		}
		// Try master branch
		zipURL = strings.TrimSuffix(repoURL, ".git") + "/archive/refs/heads/master.zip"
		resp, err = http.Get(zipURL)
	}

	if err != nil {
		return fmt.Errorf("failed to download repo zip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("zip download failed: %d", resp.StatusCode)
	}

	// Create destination directory
	os.MkdirAll(destPath, 0755)

	// Download zip file
	zipPath := filepath.Join(destPath, "repo.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, resp.Body)
	zipFile.Close()
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded zip, extracting...\n")

	// Extract zip
	err = extractZip(zipPath, destPath)
	if err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	// Clean up zip file
	os.Remove(zipPath)

	return nil
}

func extractZip(zipPath, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// GitHub zips have a top-level folder like "repo-main/"
	// We need to strip that prefix and extract contents directly to destPath
	var topDir string

	for _, f := range r.File {
		// Determine the top-level directory from first entry
		if topDir == "" {
			parts := strings.SplitN(f.Name, "/", 2)
			if len(parts) > 0 {
				topDir = parts[0] + "/"
			}
		}

		// Strip the top-level directory prefix
		relPath := strings.TrimPrefix(f.Name, topDir)
		if relPath == "" {
			continue
		}

		targetPath := filepath.Join(destPath, relPath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		// Create parent directories
		os.MkdirAll(filepath.Dir(targetPath), 0755)

		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetRepoPath(githubURL, githubRepo, githubOrg, localPath string) (string, error) {
	reposDir := ".github_repos"
	os.MkdirAll(reposDir, 0755)

	if localPath != "" {
		if _, err := os.Stat(localPath); err == nil {
			return localPath, nil
		}
		return "", fmt.Errorf("local path not found: %s", localPath)
	}

	if githubURL != "" {
		destPath := filepath.Join(reposDir, extractRepoName(githubURL))
		// Reuse if already cloned
		if hasFiles(destPath) {
			fmt.Printf("Using cached repo: %s\n", destPath)
			return destPath, nil
		}
		os.RemoveAll(destPath)
		return destPath, CloneGitHubRepo(githubURL, destPath)
	}

	if githubRepo != "" {
		repoURL := fmt.Sprintf("https://github.com/%s.git", githubRepo)
		destPath := filepath.Join(reposDir, strings.Split(githubRepo, "/")[1])
		// Reuse if already cloned
		if hasFiles(destPath) {
			fmt.Printf("Using cached repo: %s\n", destPath)
			return destPath, nil
		}
		os.RemoveAll(destPath)
		return destPath, CloneGitHubRepo(repoURL, destPath)
	}

	if githubOrg != "" {
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", githubOrg, githubOrg)
		destPath := filepath.Join(reposDir, githubOrg)
		if hasFiles(destPath) {
			fmt.Printf("Using cached repo: %s\n", destPath)
			return destPath, nil
		}
		os.RemoveAll(destPath)
		return destPath, CloneGitHubRepo(repoURL, destPath)
	}

	return "", fmt.Errorf("no repo source provided")
}

// hasFiles checks if a directory has actual source files (not just .git or repo.zip)
func hasFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, e := range entries {
		name := e.Name()
		if name != ".git" && name != "repo.zip" {
			return true
		}
	}
	return false
}

func extractRepoName(url string) string {
	parts := strings.Split(strings.TrimSuffix(url, ".git"), "/")
	return parts[len(parts)-1]
}
