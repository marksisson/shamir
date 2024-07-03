package main

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "strings"

  "github.com/hashicorp/vault/shamir"
)

func main() {
  if len(os.Args) < 2 {
    printUsage()
    os.Exit(1)
  }

  switch os.Args[1] {
  case "split":
    splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
    parts := splitCmd.Int("parts", 2, "Number of parts to split the secret into")
    threshold := splitCmd.Int("threshold", 2, "Threshold number of parts required to reconstruct the secret")

    splitCmd.Usage = func() {
      fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "  %s split [options] <file|string>\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "Options:\n")
      splitCmd.PrintDefaults()
    }

    splitCmd.Parse(os.Args[2:])

    if splitCmd.NArg() == 1 {
      // One argument provided, treat as file or string
      handleSplit(*parts, *threshold, splitCmd.Arg(0))
    } else if splitCmd.NArg() == 0 {
      // No arguments provided, read from stdin
      handleSplit(*parts, *threshold, "")
    } else {
      splitCmd.Usage()
      os.Exit(1)
    }

  case "combine":
    combineCmd := flag.NewFlagSet("combine", flag.ExitOnError)

    combineCmd.Usage = func() {
      fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "  %s combine <file1> <file2> ...\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "  %s combine  # Reads from stdin\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "Options:\n")
      combineCmd.PrintDefaults()
    }

    combineCmd.Parse(os.Args[2:])

    if combineCmd.NArg() == 0 {
      combineCmd.Usage()
      os.Exit(1)
    }

    handleCombine(combineCmd.Args())

  default:
    printUsage()
    os.Exit(1)
  }
}

func handleSplit(parts, threshold int, input string) {
  var secret []byte
  var err error

  if input == "" {
    secret, err = readFromStdin()
    if err != nil {
      log.Fatalf("Error reading secret from stdin: %v", err)
    }
  } else if fileExists(input) {
    secret, err = os.ReadFile(input)
    if err != nil {
      log.Fatalf("Error reading secret from file: %v", err)
    }
  } else {
    secret = []byte(input)
  }

  shares, err := shamir.Split(secret, parts, threshold)
  if err != nil {
    log.Fatalf("Error splitting secret: %v", err)
  }

  for i, share := range shares {
    shareFile := fmt.Sprintf("share_%d.txt", i+1)
    err := os.WriteFile(shareFile, share, 0644)
    if err != nil {
      log.Fatalf("Error writing share file: %v", err)
    }
    fmt.Printf("Share %d written to %s\n", i+1, shareFile)
  }
}

func handleCombine(shareFiles []string) {
  var shares [][]byte
  for _, shareFile := range shareFiles {
    share, err := os.ReadFile(strings.TrimSpace(shareFile))
    if err != nil {
      log.Fatalf("Error reading share file %s: %v", shareFile, err)
    }
    shares = append(shares, share)
  }

  secret, err := shamir.Combine(shares)
  if err != nil {
    log.Fatalf("Error combining shares: %v", err)
  }

  fmt.Printf("Recovered secret: %s\n", secret)
}

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }
  return !info.IsDir()
}

func readFromStdin() ([]byte, error) {
  fmt.Println("Enter secret (end with EOF):")
  reader := bufio.NewReader(os.Stdin)
  return reader.ReadBytes('\n')
}

func readLinesFromStdin() []string {
  fmt.Println("Enter share file paths (end with EOF):")
  scanner := bufio.NewScanner(os.Stdin)
  var lines []string
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    log.Fatalf("Error reading stdin: %v", err)
  }
  return lines
}

func printUsage() {
  fmt.Fprintf(os.Stderr, "Usage:\n")
  fmt.Fprintf(os.Stderr, "  shamir split [-parts <parts>] [-threshold <threshold>] <file|string>\n")
  fmt.Fprintf(os.Stderr, "  shamir split [-parts <parts>] [-threshold <threshold>]  # Reads from stdin\n")
  fmt.Fprintf(os.Stderr, "  shamir combine <file1> <file2> ...\n")
  fmt.Fprintf(os.Stderr, "  shamir combine  # Reads from stdin\n")
  fmt.Fprintf(os.Stderr, "\nCommands:\n")
  fmt.Fprintf(os.Stderr, "  split     Split a secret into parts\n")
  fmt.Fprintf(os.Stderr, "  combine   Combine parts to recover a secret\n")
}
