package application

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	goModFilename = "go.mod"
	vendorDir     = "vendor"
)

type application struct {
	drawer OutputDrawer
}

type OutputDrawer interface {
	Draw(parameters []ProjectParameters) (string, error)
}

func NewApplication(drawer OutputDrawer) *application {
	return &application{
		drawer: drawer,
	}
}

func (a *application) Execute(input Input) error {
	repositories, err := findGoModules(input.RootDirectory, input.IgnoreVendors)
	if err != nil {
		return err
	}

	output, err := a.drawer.Draw(repositories)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil
}

func findGoModules(root string, ignoreVendors bool) ([]ProjectParameters, error) {
	var repositories []ProjectParameters
	err := filepath.WalkDir(
		root,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if ignoreVendors && strings.Contains(
				path,
				fmt.Sprintf("%v%s%v", filepath.Separator, vendorDir, filepath.Separator),
			) {
				return nil
			}

			if d.Name() == goModFilename {
				goVersion, err := extractGoVersion(path)
				if err != nil {
					return err
				}

				repositoryName, repositoryUrl := findGitRepo(path)
				if repositoryName != "" {
					repositories = append(
						repositories,
						ProjectParameters{
							GoVersion:   goVersion,
							ProjectName: repositoryName,
							ProjectUrl:  repositoryUrl,
						},
					)
				}
			}

			return nil
		},
	)

	return repositories, err
}

func extractGoVersion(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	const goModVersionPrefix = "go "
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, goModVersionPrefix) {
			return strings.TrimPrefix(line, goModVersionPrefix), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func findGitRepo(modPath string) (string, string) {
	dir := filepath.Dir(modPath)

	for {
		gitConfigPath := filepath.Join(dir, ".git", "config")
		if _, err := os.Stat(gitConfigPath); err == nil {
			return extractGitRemote(gitConfigPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return "", ""
}

func extractGitRemote(configPath string) (string, string) {
	file, err := os.Open(configPath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRemoteOrigin := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "[remote \"origin\"]") {
			inRemoteOrigin = true
			continue
		}
		if inRemoteOrigin && strings.HasPrefix(line, "url = ") {
			url := strings.TrimPrefix(line, "url = ")
			repoName := extractRepositoryName(url)
			return repoName, url
		}
	}

	return "", ""
}

func extractRepositoryName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
