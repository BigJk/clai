package templating

import (
	"testing"

	"github.com/bigjk/clai/ai"
	"github.com/stretchr/testify/assert"
)

func TestParseTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		want     []ai.Message
	}{
		{
			name: "simple template",
			template: `# CLAI::SYSTEM
You are a helpful assistant.

# CLAI::USER
Hello, how are you?`,
			want: []ai.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "Hello, how are you?"},
			},
		},
		{
			name: "empty lines between roles",
			template: `# CLAI::SYSTEM
You are a helpful assistant.


# CLAI::USER

Hello, how are you?

`,
			want: []ai.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "Hello, how are you?"},
			},
		},
		{
			name: "multiple roles with content",
			template: `# CLAI::SYSTEM
You are a data generator.

# CLAI::USER
First user message

# CLAI::ASSISTANT
Assistant response

# CLAI::USER
Second user message`,
			want: []ai.Message{
				{Role: "system", Content: "You are a data generator."},
				{Role: "user", Content: "First user message"},
				{Role: "assistant", Content: "Assistant response"},
				{Role: "user", Content: "Second user message"},
			},
		},
		{
			name: "content with multiple lines",
			template: `# CLAI::SYSTEM
Line 1
Line 2
Line 3

# CLAI::USER
Multi
Line
Message`,
			want: []ai.Message{
				{Role: "system", Content: "Line 1\nLine 2\nLine 3"},
				{Role: "user", Content: "Multi\nLine\nMessage"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTemplate(tt.template)
			assert.Equal(t, tt.want, got)
		})
	}
}
