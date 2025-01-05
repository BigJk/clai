package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bigjk/clai/ai"
	"github.com/bigjk/clai/executor"
	"github.com/bigjk/clai/templating"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	URL    string
	APIKey string
	Model  string
}

var Version = "dev"

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	viper.SetConfigName(".clairc")
	viper.SetConfigType("yaml")

	// Add config paths in order of increasing priority
	viper.AddConfigPath(home) // $HOME/.clairc
	viper.AddConfigPath(".")  // ./.clairc (current directory)

	// Set required fields
	viper.SetDefault("url", "")
	viper.SetDefault("apikey", "")
	viper.SetDefault("model", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func createConfigCmd() *cobra.Command {
	var (
		useOpenAI     bool
		useOpenRouter bool
	)

	cmd := &cobra.Command{
		Use:   "create-config",
		Short: "Create a new .clairc file in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if useOpenAI && useOpenRouter {
				return fmt.Errorf("cannot use both --openai and --open_router flags")
			}

			config := map[string]string{
				"url":    "",
				"apikey": "",
				"model":  "",
			}

			if useOpenAI {
				config["url"] = "https://api.openai.com/v1/chat/completions"
			} else if useOpenRouter {
				config["url"] = "https://openrouter.ai/api/v1/chat/completions"
			}

			viper.SetConfigFile(".clairc")
			viper.SetConfigType("yaml")
			for k, v := range config {
				viper.Set(k, v)
			}

			return viper.WriteConfig()
		},
	}

	cmd.Flags().BoolVar(&useOpenAI, "openai", false, "Use OpenAI as the provider")
	cmd.Flags().BoolVar(&useOpenRouter, "open_router", false, "Use OpenRouter as the provider")

	return cmd
}

func runCmd() *cobra.Command {
	var (
		workingDir string
		outFile    string
		dryRun     bool
	)

	cmd := &cobra.Command{
		Use:   "run [file] [input...]",
		Short: "Run a file with the given input",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			input := strings.Join(args[1:], " ")

			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			messages := templating.ParseTemplate(string(content))

			finalMessages, err := executor.Execute(messages, input, workingDir)
			if err != nil {
				return fmt.Errorf("error executing command: %w", err)
			}

			var result string
			if dryRun {
				// Format messages for preview
				result = "Messages that would be sent to API:\n\n"
				for i, msg := range finalMessages {
					result += fmt.Sprintf("Message %d:\n", i+1)
					result += fmt.Sprintf("Role: %s\n", msg.Role)
					result += fmt.Sprintf("Content:\n%s\n\n", msg.Content)
				}
			} else {
				res, err := ai.NewClient(
					ai.WithAPIKey(viper.GetString("apikey")),
					ai.WithModel(viper.GetString("model")),
					ai.WithURL(viper.GetString("url")),
				).Do(finalMessages)
				if err != nil {
					return fmt.Errorf("error getting response: %w", err)
				}
				result = res
			}

			if outFile != "" {
				if err := os.WriteFile(outFile, []byte(result), 0644); err != nil {
					return fmt.Errorf("error writing result file: %w", err)
				}
			} else {
				fmt.Println(result)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&workingDir, "working_dir", "./", "Working directory for the command")
	cmd.Flags().StringVar(&outFile, "out", "", "Output file path (if not specified, prints to stdout)")
	cmd.Flags().BoolVar(&dryRun, "dry", false, "Preview messages without sending to API")
	return cmd
}

func runMultipleCmd() *cobra.Command {
	var (
		workingDir string
		outDir     string
		numRuns    int
		dryRun     bool
	)

	cmd := &cobra.Command{
		Use:   "run_multiple [file] [input...]",
		Short: "Run a file multiple times with the given input and save results",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			input := strings.Join(args[1:], " ")

			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			messages := templating.ParseTemplate(string(content))
			client := ai.NewClient(
				ai.WithAPIKey(viper.GetString("apikey")),
				ai.WithModel(viper.GetString("model")),
				ai.WithURL(viper.GetString("url")),
			)

			var errors []error
			wg := &sync.WaitGroup{}
			wg.Add(numRuns)

			for i := 0; i < numRuns; i++ {
				go func(i int) {
					defer wg.Done()

					finalMessages, err := executor.Execute(messages, input, workingDir)
					if err != nil {
						errors = append(errors, fmt.Errorf("error executing command: %w", err))
						return
					}

					var result string
					if dryRun {
						// Format messages for preview
						result = fmt.Sprintf("Run %d - Messages that would be sent to API:\n\n", i+1)
						for j, msg := range finalMessages {
							result += fmt.Sprintf("Message %d:\n", j+1)
							result += fmt.Sprintf("Role: %s\n", msg.Role)
							result += fmt.Sprintf("Content:\n%s\n\n", msg.Content)
						}
					} else {
						res, err := client.Do(finalMessages)
						if err != nil {
							errors = append(errors, fmt.Errorf("error getting response: %w", err))
							return
						}
						result = res
					}

					outFile := filepath.Join(outDir, fmt.Sprintf("res_%d.md", i+1))
					if err := os.WriteFile(outFile, []byte(result), 0644); err != nil {
						errors = append(errors, fmt.Errorf("error writing result file: %w", err))
						return
					}
				}(i)
			}

			wg.Wait()

			if len(errors) > 0 {
				return fmt.Errorf("errors occurred: %v", errors)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&workingDir, "working_dir", "./", "Working directory for the command")
	cmd.Flags().StringVar(&outDir, "out", "./", "Output directory for result files")
	cmd.Flags().IntVar(&numRuns, "num", 3, "Number of times to run the workflow")
	cmd.Flags().BoolVar(&dryRun, "dry", false, "Preview messages without sending to API")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("clai version %s\n", Version)
		},
	}
}

func main() {
	loadConfig()

	rootCmd := &cobra.Command{
		Use:     "clai",
		Short:   "CLAI - Command Line AI Workflow Runner",
		Version: Version,
	}

	rootCmd.AddCommand(createConfigCmd())
	rootCmd.AddCommand(runCmd())
	rootCmd.AddCommand(runMultipleCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
