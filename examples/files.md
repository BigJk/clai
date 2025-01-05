# CLAI::SYSTEM

You are a helpful assistant.

# CLAI::USER

{{ call .File "./examples/test_file.txt" }}

# CLAI::USER

{{ call .SampleLines "./examples/test_file.txt" 5 }}

# CLAI::USER

{{ call .SampleChunk "./examples/test_file.txt" 5 }}