package main

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "testing"
)

func TestSplitFromFile(t *testing.T) {
  secret := "my secret"
  tmpfile, err := os.CreateTemp("", "secret.txt")
  if err != nil {
    t.Fatal(err)
  }
  defer os.Remove(tmpfile.Name())

  if _, err := tmpfile.WriteString(secret); err != nil {
    t.Fatal(err)
  }
  if err := tmpfile.Close(); err != nil {
    t.Fatal(err)
  }

  cmd := exec.Command("go", "run", "shamir.go", "split", "-parts", "3", "-threshold", "2", tmpfile.Name())
  var out, stderr bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &stderr
  if err := cmd.Run(); err != nil {
    t.Fatalf("split command failed: %v, stderr: %v", err, stderr.String())
  }

  for i := 1; i <= 3; i++ {
    shareFile := fmt.Sprintf("share_%d.txt", i)
    if _, err := os.Stat(shareFile); os.IsNotExist(err) {
      t.Fatalf("expected share file %s not found", shareFile)
    }
    defer os.Remove(shareFile)
  }
}

func TestSplitFromString(t *testing.T) {
  cmd := exec.Command("go", "run", "shamir.go", "split", "-parts", "3", "-threshold", "2", "my secret string")
  var out, stderr bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &stderr
  if err := cmd.Run(); err != nil {
    t.Fatalf("split command failed: %v, stderr: %v", err, stderr.String())
  }

  for i := 1; i <= 3; i++ {
    shareFile := fmt.Sprintf("share_%d.txt", i)
    if _, err := os.Stat(shareFile); os.IsNotExist(err) {
      t.Fatalf("expected share file %s not found", shareFile)
    }
    defer os.Remove(shareFile)
  }
}

func TestSplitFromStdin(t *testing.T) {
  cmd := exec.Command("go", "run", "shamir.go", "split", "-parts", "3", "-threshold", "2")
  cmd.Stdin = strings.NewReader("my secret from stdin\n")
  var out, stderr bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &stderr
  if err := cmd.Run(); err != nil {
    t.Fatalf("split command failed: %v, stderr: %v", err, stderr.String())
  }

  for i := 1; i <= 3; i++ {
    shareFile := fmt.Sprintf("share_%d.txt", i)
    if _, err := os.Stat(shareFile); os.IsNotExist(err) {
      t.Fatalf("expected share file %s not found", shareFile)
    }
    defer os.Remove(shareFile)
  }
}

func TestCombineFromFiles(t *testing.T) {
  // First, split a secret to generate shares
  splitCmd := exec.Command("go", "run", "shamir.go", "split", "-parts", "3", "-threshold", "2", "my secret string")
  var splitOut, splitErr bytes.Buffer
  splitCmd.Stdout = &splitOut
  splitCmd.Stderr = &splitErr
  if err := splitCmd.Run(); err != nil {
    t.Fatalf("split command failed: %v, stderr: %v", err, splitErr.String())
  }

  // Then, combine the generated shares
  combineCmd := exec.Command("go", "run", "shamir.go", "combine", "share_1.txt", "share_2.txt")
  var combineOut, combineErr bytes.Buffer
  combineCmd.Stdout = &combineOut
  combineCmd.Stderr = &combineErr
  if err := combineCmd.Run(); err != nil {
    t.Fatalf("combine command failed: %v, stderr: %v", err, combineErr.String())
  }

  expected := "Recovered secret: my secret string\n"
  if combineOut.String() != expected {
    t.Fatalf("expected %q, got %q", expected, combineOut.String())
  }

  for i := 1; i <= 3; i++ {
    shareFile := fmt.Sprintf("share_%d.txt", i)
    defer os.Remove(shareFile)
  }
}

func TestCombineFromWildcard(t *testing.T) {
  // First, split a secret to generate shares
  splitCmd := exec.Command("go", "run", "shamir.go", "split", "-parts", "3", "-threshold", "2", "my secret string")
  var splitOut, splitErr bytes.Buffer
  splitCmd.Stdout = &splitOut
  splitCmd.Stderr = &splitErr
  if err := splitCmd.Run(); err != nil {
    t.Fatalf("split command failed: %v, stderr: %v", err, splitErr.String())
  }

  // Expand wildcard pattern
  files, err := filepath.Glob("share_*.txt")
  if err != nil {
    t.Fatalf("failed to expand wildcard: %v", err)
  }
  if len(files) == 0 {
    t.Fatalf("no files matched the wildcard pattern")
  }

  // Then, combine the generated shares using the expanded list
  args := append([]string{"run", "shamir.go", "combine"}, files...)
	combineCmd := exec.Command("go", args...)
  var combineOut, combineErr bytes.Buffer
  combineCmd.Stdout = &combineOut
  combineCmd.Stderr = &combineErr
  if err := combineCmd.Run(); err != nil {
    t.Fatalf("combine command failed: %v, stderr: %v", err, combineErr.String())
  }

  expected := "Recovered secret: my secret string\n"
  if combineOut.String() != expected {
    t.Fatalf("expected %q, got %q", expected, combineOut.String())
  }

  for _, file := range files {
    defer os.Remove(file)
  }
}
