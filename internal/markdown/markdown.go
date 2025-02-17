package markdown

import (
	"fmt"
	"strings"

	"github.com/vehsamrak/go-version-parser/internal/application"
)

type markdownDrawer struct{}

func NewMarkdownDrawer() *markdownDrawer {
	return &markdownDrawer{}
}

func (m *markdownDrawer) Draw(projects []application.ProjectParameters) (string, error) {
	if len(projects) == 0 {
		return "", fmt.Errorf("no projects found")
	}

	var sb strings.Builder
	sb.WriteString("| Repository | Go version |\n")
	sb.WriteString("|------------|------------|\n")

	for _, project := range projects {
		sb.WriteString(fmt.Sprintf(
			"| [%s](%s) | %s |\n",
			project.ProjectName,
			project.ProjectUrl,
			project.GoVersion,
		))
	}

	return sb.String(), nil
}
