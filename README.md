# DirSleuth - Directory Enumeration Tool

## Overview
DirSleuth is a fast and efficient directory enumeration tool written in Go. It scans target domains to discover hidden directories and files using a customizable wordlist. Designed for penetration testers and developers, it is lightweight, user-friendly, and highly performant.

## Features
- Multithreaded scanning for improved performance.
- Support for customizable wordlists.
- Configurable options for request rate and timeout.
- Clean and modular codebase for easy contributions.

## Installation

### Prerequisites
- [Go Programming Language](https://go.dev/dl/)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/Coding-for-Weeks/dirsleuth.git
   cd dirsleuth
   ```
2. Build the binary:
   ```bash
   make build
   ```
3. Verify the binary is created:
   ```bash
   ls -l DirSleuth
   ```

## Usage

### Basic Command
Run the tool with a target domain and a wordlist:
```bash
./DirSleuth -d example.com -w /path/to/wordlist.txt
```

### Options
- `-d`: The target domain to scan.
- `-w`: Path to the wordlist file.
- `-t`: Number of threads to use for scanning (default: 10).

### Examples
1. Scan with default settings:
   ```bash
    ./DirSleuth -d example.com -w wordlist.txt
   ```
2. Customize threads and timeout:
   ```bash
    ./DirSleuth -d example.com -w wordlist.txt -t 20 -timeout 10
   ```

## Development

### Running Tests
Run tests to ensure code quality:
```bash
make test
```
Run specific tests or benchmarks:
```bash
make test TEST_FLAGS="-run TestSpecificFunction -v"
```

### Building and Running
Build and run the project with a custom configuration:
```bash
make run CONFIG=config.json
```

### Cleaning Up
Remove all build artifacts and temporary files:
```bash
make clean
```

## Contribution
We welcome contributions! If you want to improve DirSleuth:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with detailed changes.

### Suggested Enhancements
- Additional scanning modes.
- Improved error handling.
- Expanded output formats (e.g., JSON, XML).

## License
DirSleuth is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it under the terms of the license.

---
For questions, suggestions, or issues, please open an [issue](https://github.com/Coding-for-Weeks/dirsleuth/issues) or contact us directly.

