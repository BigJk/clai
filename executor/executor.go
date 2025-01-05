package executor

import (
	"encoding/json"
	"path/filepath"

	"github.com/bigjk/clai/ai"
	"github.com/bigjk/clai/templating"
)

func Execute(messages []ai.Message, userInput string, rootDir string) ([]ai.Message, error) {
	var data map[string]any

	err := json.Unmarshal([]byte(userInput), &data)
	if err != nil {
		data = map[string]any{"Input": userInput}
	}

	data["SampleFiles"] = func(folder string, count int, meta bool) string {
		return SampleFiles(filepath.Join(rootDir, folder), count, meta)
	}
	data["SampleLines"] = func(file string, count int) string {
		return SampleLines(filepath.Join(rootDir, file), count)
	}
	data["File"] = func(file string) string {
		return File(filepath.Join(rootDir, file))
	}
	data["SampleChunk"] = func(file string, count int) string {
		return SampleChunk(filepath.Join(rootDir, file), count)
	}
	data["SampleLines"] = func(file string, count int) string {
		return SampleLines(filepath.Join(rootDir, file), count)
	}
	data["RunCommand"] = func(command string, args ...string) string {
		if command[0] == '.' {
			command = filepath.Join(rootDir, command[1:])
		}
		return RunCommand(command, args...)
	}

	var newMessages []ai.Message
	for i := range messages {
		res, err := templating.ExecuteTemplate(messages[i].Content, data)
		if err != nil {
			return nil, err
		}
		newMessages = append(newMessages, ai.Message{
			Role:    messages[i].Role,
			Content: res,
		})
	}

	return newMessages, nil
}
