package executor

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

		content, err := ioutil.ReadFile(filepath.Join(folder, file.Name()))
		if err != nil {
			result.WriteString(fmt.Sprintf("Error reading file %s: %v\n", file.Name(), err))
			continue
		}

		if meta {
			result.WriteString(fmt.Sprintf("====== File: %s\n", file.Name()))
		}
		result.Write(content)
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

	lines := strings.Split(string(content), "\n")
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
	return string(content)
}

// SampleChunk reads a file and returns a random chunk with lines count
func SampleChunk(file string, count int) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
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
