# CL(A)I

CLAI is a command line interface to quickly create and run workflows against LLM with helper to make it easier to insert data into the workflow. Data can come from the user, commands or files. There are utilities to sample files, lines and chunks from files, to quickly build examples.

This can be helpful for creating workflows that generate new data based on existing data files. The templating engine is based on [Go's html/template](https://pkg.go.dev/html/template). By providing the templating methods it's easy to insert up-to-date data into the workflow.

## Background

I extensifly use LLMs to brainstorm ideas, especially in the context of TTRPGs like Dungeons & Dragons. I wanted to make it easier to create and run workflows against LLMs with up-to-date examples from my obsidian vault.

## Installation

### Quick Install (Linux and macOS)

This will install the latest version of `clai` to `/usr/local/bin`

```bash
curl -fsSL https://raw.githubusercontent.com/bigjk/clai/main/install.sh | bash
```

### Manual Installation

#### Go

```bash
go install github.com/bigjk/clai@latest
```

#### Binary Release
You can download the pre-built binary for your platform from the [releases page](https://github.com/bigjk/clai/releases).

## Config File

CLAI uses a YAML config file to store the API key, model, and url of the OpenAI compatible API.
The config file should be named `.clairc` and be placed in the current directory or in the `$HOME` directory. Current directory takes precedence over `$HOME`.

```yaml
url: https://api.openai.com/v1/chat/completions
apikey: YOUR_OPENAI_API_KEY
model: gpt-4o-mini
```

## CLI Usage

CLAI provides several commands to help you manage and run your workflows:

### Configuration

```bash
# Create a new config file in the current directory
clai create-config --openai      # Configure for OpenAI
clai create-config --open_router # Configure for OpenRouter
```

### Running Workflows

#### Single Run
```bash
clai run [workflow_file] [input...]
  --working_dir string   Working directory for the command (default "./")
  --out string          Output file path (if not specified, prints to stdout)
  --dry                 Preview messages without sending to API
```

Example:
```bash
# Print to stdout
clai run ./workflow.md "Generate a creative story about a space adventure"

# Save to file
clai run --out "./result.md" ./workflow.md "Generate a creative story about a space adventure"

# Preview messages without API call
clai run --dry ./workflow.md "Generate a creative story about a space adventure"
```

#### Multiple Parallel Runs
```bash
clai run_multiple [workflow_file] [input...]
  --working_dir string   Working directory for the command (default "./")
  --out string          Output directory for result files (default "./")
  --num int            Number of times to run the workflow (default 3)
  --dry                Preview messages without sending to API
```

Example:
```bash
# Run the workflow 5 times in parallel and save results as res_1.md through res_5.md
clai run_multiple --num 5 --out "./results" ./workflow.md "Generate different variations of a product description"

# Preview messages for 5 runs without API calls
clai run_multiple --dry --num 5 --out "./results" ./workflow.md "Generate different variations of a product description"
```

### Template Functions

In your workflow files, you can use several helper functions:

- `{{ .Input }}`: Insert the user's input (when using plain text input)
- `{{ .CustomInput }}`, `{{ .Field1 }}`, etc.: Access JSON fields (when using JSON input)
- `{{ call .SampleFiles "path" n false }}`: Sample n random files from the specified path
- `{{ call .SampleFiles "path" n true }}`: Sample n random files with their filenames as headers
- `{{ call .SampleLines "file" n }}`: Sample n random lines from the specified file
- `{{ call .File "path" }}`: Read and return the entire contents of a file
- `{{ call .SampleChunk "file" n }}`: Read a random chunk of n consecutive lines from a file
- `{{ call .RunCommand "cmd" "arg1" "arg2" }}`: Execute a shell command and return its output

### Input Types

CLAI supports both plain text and JSON input formats:

#### Plain Text Input
```markdown
# CLAI::SYSTEM
You are a helpful assistant.

# CLAI::USER
{{ .Input }}
```

```bash
clai run ./workflow.md "Hello, how are you?"
```

#### JSON Input
```markdown
# CLAI::SYSTEM
You are a helpful assistant.

# CLAI::USER
The user's custom input is: {{ .CustomInput }}
Additional field: {{ .Field1 }}
```

```bash
# Access JSON fields in the template
clai run ./workflow.md '{ "CustomInput": "Hello!", "Field1": "Extra data" }'

# Preview the processed template
clai run --dry ./workflow.md '{ "CustomInput": "Hello!", "Field1": "Extra data" }'
```

The JSON input is parsed and its fields become available in the template using dot notation. This is useful when you need to pass structured data to your workflow.

## Example Workflow File

### TTRPG Example

Imagine you have a directory with markdown files of creatures for a tabletop roleplaying game. You want to generate new creatures for the game.

```markdown
# CLAI::SYSTEM

You are a helpful assistant that generates new monsters for a tabletop roleplaying game. You will generate a new monster based on the given input.

# CLAI::USER

Here are some examples of monsters:

{{ call .SampleFiles "./monsters/" 5 true }}

# CLAI::ASSISTANT

Thank you for the examples! Now tell me about the monster you want to generate.

# CLAI::USER

{{ .Input }}
```

```bash
# Single run print to stdout
clai run ./monsters.md "A dragon with a fire breath attack"

# Single run save to file
clai run --out "./result.md" ./monsters.md "A dragon with a fire breath attack"

# Multiple runs
clai run_multiple --num 5 --out "./results" ./monsters.md "A dragon with a fire breath attack"