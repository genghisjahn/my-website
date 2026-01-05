---
slug: go-error-handling
title: Go Error Handling Pattern
date: 2026-01-05
author:
  name: Jon Wear
tags:
  - name: go
    slug: go
  - name: tips
    slug: tips
draft: false
---

A common pattern for handling errors in Go with context wrapping:

```go
func doSomething() error {
    result, err := someOperation()
    if err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

Using `%w` allows error unwrapping with `errors.Is()` and `errors.As()`.

## Multiple Error Checks

```go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open file: %w", err)
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
        return fmt.Errorf("read file: %w", err)
    }

    // use data...
    return nil
}
```
