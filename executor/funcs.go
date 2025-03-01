package executor

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SampleFiles reads count random files from the folder and appends them as a string.
// If meta is set the file name is included
func SampleFiles(folder string, count int, meta bool) string {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	if count > len(files) {
		count = len(files)
	}

	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	var result strings.Builder
	for i := 0; i < count; i++ {
		file := files[i]
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(folder, file.Name()))
		if err != nil {
			result.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file.Name(), err))
			continue
		}

		if meta {
			result.WriteString(fmt.Sprintf("====== File: %s\n", file.Name()))
		}
		result.WriteString(RemoveFrontmatter(file.Name(), string(content)))
		result.WriteString("\n\n")
	}

	return result.String()
}

// SampleFilesDeep reads count random files from the folder and appends them as a string.
// If meta is set the file name is included. This is a recursive function.
func SampleFilesDeep(folder string, count int, meta bool) string {
	var possibleFiles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			possibleFiles = append(possibleFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	rand.Shuffle(len(possibleFiles), func(i, j int) {
		possibleFiles[i], possibleFiles[j] = possibleFiles[j], possibleFiles[i]
	})

	var result strings.Builder
	for i := 0; i < count; i++ {
		file := possibleFiles[i]
		content, err := os.ReadFile(file)
		if err != nil {
			result.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file, err))
			continue
		}

		if meta {
			result.WriteString(fmt.Sprintf("====== File: %s\n", file))
		}
		result.WriteString(RemoveFrontmatter(file, string(content)))
		result.WriteString("\n\n")
	}

	return result.String()
}

// SampleLines reads count random lines from the file and appends them as a string
func SampleLines(file string, count int) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(RemoveFrontmatter(file, string(content)), "\n")
	if count > len(lines) {
		count = len(lines)
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	var result strings.Builder
	for i := 0; i < count; i++ {
		result.WriteString(lines[i])
		result.WriteString("\n")
	}

	return result.String()
}

// File reads a file and returns its content
func File(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return RemoveFrontmatter(file, string(content))
}

// SampleChunk reads a file and returns a random chunk with lines count
func SampleChunk(file string, count int) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(RemoveFrontmatter(file, string(content)), "\n")
	if count > len(lines) {
		count = len(lines)
	}

	start := rand.Intn(len(lines) - count)
	end := start + count
	return strings.Join(lines[start:end], "\n")
}

// RunCommand runs a command and returns its output
func RunCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// SampleFilesPattern reads count random files from the folder whose content matches the pattern and appends them as a string.
// If meta is set the file name is included
func SampleFilesPattern(folder string, pattern string, count int, meta bool) string {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	var matchingFiles []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(folder, file.Name()))
		if err != nil {
			continue
		}

		matched, err := regexp.MatchString(pattern, string(content))
		if err != nil {
			continue
		}

		if matched {
			matchingFiles = append(matchingFiles, file)
		}
	}

	if count > len(matchingFiles) {
		count = len(matchingFiles)
	}

	rand.Shuffle(len(matchingFiles), func(i, j int) {
		matchingFiles[i], matchingFiles[j] = matchingFiles[j], matchingFiles[i]
	})

	var result strings.Builder
	for i := 0; i < count; i++ {
		file := matchingFiles[i]
		content, err := os.ReadFile(filepath.Join(folder, file.Name()))
		if err != nil {
			result.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file.Name(), err))
			continue
		}

		if meta {
			result.WriteString(fmt.Sprintf("====== File: %s\n", file.Name()))
		}
		result.WriteString(RemoveFrontmatter(file.Name(), string(content)))
		result.WriteString("\n\n")
	}

	return result.String()
}

// SampleFilesPatternDeep reads count random files from the folder and its subdirectories whose content matches the pattern
// and appends them as a string. If meta is set the file name is included.
func SampleFilesPatternDeep(folder string, pattern string, count int, meta bool) string {
	var matchingFiles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			matched, err := regexp.MatchString(pattern, string(content))
			if err != nil {
				return nil
			}

			if matched {
				matchingFiles = append(matchingFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	if count > len(matchingFiles) {
		count = len(matchingFiles)
	}

	rand.Shuffle(len(matchingFiles), func(i, j int) {
		matchingFiles[i], matchingFiles[j] = matchingFiles[j], matchingFiles[i]
	})

	var result strings.Builder
	for i := 0; i < count; i++ {
		file := matchingFiles[i]
		content, err := os.ReadFile(file)
		if err != nil {
			result.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file, err))
			continue
		}

		if meta {
			result.WriteString(fmt.Sprintf("====== File: %s\n", file))
		}
		result.WriteString(RemoveFrontmatter(file, string(content)))
		result.WriteString("\n\n")
	}

	return result.String()
}
