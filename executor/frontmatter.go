package executor

import (
	"path/filepath"
	"strings"
)

// RemoveFrontmatter removes YAML frontmatter from markdown content.
// Frontmatter must be enclosed between "---" lines at the start of the content.
func RemoveFrontmatter(file string, content string) string {
	// Only process markdown files
	if !strings.HasSuffix(filepath.Ext(file), ".md") {
		return content
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 2 || !strings.HasPrefix(lines[0], "---") {
		return content
	}

	// Find the closing frontmatter delimiter
	for i := 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "---") {
			// Return content after frontmatter
			return strings.Join(lines[i+1:], "\n")
		}
	}

	// If we didn't find a closing delimiter, return original content
	return content
}
