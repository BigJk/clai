package templating

import (
	"strings"

	"github.com/bigjk/clai/ai"
)

// ParseTemplate takes a template string and returns a slice of ai.Message
func ParseTemplate(template string) []ai.Message {
	var messages []ai.Message
	var currentRole string
	var currentContent []string

	lines := strings.Split(template, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "# CLAI::") {
			// If we have a previous role and content, add it to messages
			if currentRole != "" && len(currentContent) > 0 {
				messages = append(messages, ai.Message{
					Role:    strings.ToLower(currentRole),
					Content: strings.TrimSpace(strings.Join(currentContent, "\n")),
				})
				currentContent = nil
			}
			// Extract new role
			currentRole = strings.TrimPrefix(trimmedLine, "# CLAI::")
		} else {
			// Add line to current content if we have a role
			if currentRole != "" {
				currentContent = append(currentContent, line)
			}
		}
	}

	// Add the last message if there is one
	if currentRole != "" && len(currentContent) > 0 {
		messages = append(messages, ai.Message{
			Role:    strings.ToLower(currentRole),
			Content: strings.TrimSpace(strings.Join(currentContent, "\n")),
		})
	}

	return messages
}
