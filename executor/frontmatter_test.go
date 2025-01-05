package executor

import "testing"

func TestRemoveFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		content  string
		expected string
	}{
		{
			name:     "non markdown file",
			file:     "test.txt",
			content:  "---\ntitle: Test\n---\nHello",
			expected: "---\ntitle: Test\n---\nHello",
		},
		{
			name:     "markdown without frontmatter",
			file:     "test.md",
			content:  "# Hello\nWorld",
			expected: "# Hello\nWorld",
		},
		{
			name:     "markdown with frontmatter",
			file:     "test.md",
			content:  "---\ntitle: Test\n---\n# Hello\nWorld",
			expected: "# Hello\nWorld",
		},
		{
			name:     "markdown with incomplete frontmatter",
			file:     "test.md",
			content:  "---\ntitle: Test\n# Hello\nWorld",
			expected: "---\ntitle: Test\n# Hello\nWorld",
		},
		{
			name:     "markdown with empty frontmatter",
			file:     "test.md",
			content:  "---\n---\n# Hello\nWorld",
			expected: "# Hello\nWorld",
		},
		{
			name:     "markdown with complex frontmatter",
			file:     "test.md",
			content:  "---\ntitle: Test\ndescription: |\n  Multiple\n  Lines\ntags:\n  - one\n  - two\n---\n# Hello\nWorld",
			expected: "# Hello\nWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveFrontmatter(tt.file, tt.content)
			if result != tt.expected {
				t.Errorf("RemoveFrontmatter() = %q, want %q", result, tt.expected)
			}
		})
	}
}
