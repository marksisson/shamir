# Shamir's Secret Sharing CLI Tool

A command-line tool written in Go to securely split and combine secrets using Shamir's Secret Sharing Scheme. This tool allows you to split a secret into multiple parts, where a specified threshold number of parts are required to reconstruct the original secret.

## Features

- **Split a Secret:** Divide a secret into multiple parts.
- **Combine Shares:** Reconstruct the original secret from the specified number of shares.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16+ recommended)

### Install the CLI tool

```sh
go install github.com/marksisson/shamir@latest
```

### Add GOPATH/bin to your PATH (if not already done)

```sh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Verify the installation

```sh
shamir --help
```

## Usage

### To Split a Secret

#### From a file

```sh
shamir split -parts 5 -threshold 3 secret.txt
```

#### From a string

```sh
shamir split -parts 5 -threshold 3 "my secret string"
```

#### From stdin

```sh
shamir split -parts 5 -threshold 3
```

Then type or paste the secret followed by `Ctrl+D`.

### To Combine Shares

#### From specific files

```sh
shamir combine share_1.txt share_2.txt share_3.txt
```

#### Using wildcard (expanded by shell)

```sh
shamir combine share_*.txt
```

#### From stdin

```sh
shamir combine
```

Then type or paste the share file paths followed by `Ctrl+D`.

## Development

### Running Tests

Ensure you have Go installed and your environment set up. Clone the repository and run the tests using:

```sh
git clone https://github.com/marksisson/shamir.git
cd shamir
go test -v
```

## License

This project is licensed under the terms of the [MIT license](LICENSE).
