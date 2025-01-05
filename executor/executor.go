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

	registerFunc := func(names []string, f any) {
		for _, name := range names {
			data[name] = f
		}
	}

	registerFunc([]string{"SampleFiles", "SF"}, func(folder string, count int, meta bool) string {
		return SampleFiles(filepath.Join(rootDir, folder), count, meta)
	})
	registerFunc([]string{"SampleFilesDeep", "SFD"}, func(folder string, count int, meta bool) string {
		return SampleFilesDeep(filepath.Join(rootDir, folder), count, meta)
	})
	registerFunc([]string{"SampleFilesPattern", "SFP"}, func(folder string, pattern string, count int, meta bool) string {
		return SampleFilesPattern(filepath.Join(rootDir, folder), pattern, count, meta)
	})
	registerFunc([]string{"SampleFilesPatternDeep", "SFDP"}, func(folder string, pattern string, count int, meta bool) string {
		return SampleFilesPatternDeep(filepath.Join(rootDir, folder), pattern, count, meta)
	})
	registerFunc([]string{"SampleLines", "SL"}, func(file string, count int) string {
		return SampleLines(filepath.Join(rootDir, file), count)
	})
	registerFunc([]string{"File", "F"}, func(file string) string {
		return File(filepath.Join(rootDir, file))
	})
	registerFunc([]string{"SampleChunk", "SC"}, func(file string, count int) string {
		return SampleChunk(filepath.Join(rootDir, file), count)
	})
	registerFunc([]string{"RunCommand", "RC"}, func(command string, args ...string) string {
		if command[0] == '.' {
			command = filepath.Join(rootDir, command[1:])
		}
		return RunCommand(command, args...)
	})

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
