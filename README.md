# go-promise

A simple Promise implementation for Go, inspired by JavaScript's Promise pattern. This library allows you to write asynchronous code in a more sequential and readable manner, with support for chaining, error handling, and cleanup operations.

## Features

- Create promises that resolve or reject
- Chain operations with `Then()` method
- Handle errors with `Catch()` method
- Execute cleanup code with `Finally()` method
- Generic implementation supporting any data type
- Utility function to wait for all promises to complete

## Installation

```bash
go get github.com/saatvikAwasthi/promise
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    
    "promise/pkg"
)

func main() {
    // Create a promise that resolves with a string
    pkg.NewPromise[string](func(resolve func(string), reject func(error), finally func()) {
        // Simulate an async operation
        fmt.Println("Starting async work...")
        time.Sleep(2 * time.Second)
        resolve("Operation completed successfully")
    }).Then(func(result string) {
        // Handle successful result
        fmt.Printf("Success: %s\n", result)
    }).Catch(func(err error) {
        // Handle errors
        fmt.Printf("Error: %s\n", err.Error())
    })
    
    // Wait for all promises to complete
    pkg.WaitForPromises()
}
```

### Chaining Operations

```go
pkg.NewPromise[string](func(resolve func(string), reject func(error), finally func()) {
    resolve("Initial value")
}).Then(func(result string) string {
    return result + " - transformed"
}).Then(func(result string) {
    fmt.Println(result) // Outputs: "Initial value - transformed"
})
```

### Error Handling

```go
pkg.NewPromise[int](func(resolve func(int), reject func(error), finally func()) {
    reject(fmt.Errorf("something went wrong"))
}).Then(func(result int) {
    // This won't execute
    fmt.Println(result)
}).Catch(func(err error) {
    // This will execute
    fmt.Printf("Caught error: %s\n", err.Error())
})
```

### Cleanup with Finally

```go
pkg.NewPromise[string](func(resolve func(string), reject func(error), finally func()) {
    // Do work...
    resolve("data")
    // or reject(err)
}).Finally(func() {
    // This will always execute
    fmt.Println("Cleaning up resources...")
})
```

## API Reference

### `NewPromise[T]`

Creates a new promise that can resolve with a value of type T or reject with an error.

### `Then(func(T))`

Attaches a callback that receives the resolved value.

### `Catch(func(error))`

Attaches a callback that handles errors if the promise rejects.

### `Finally(func())`

Attaches a callback that executes regardless of whether the promise resolves or rejects.

### `WaitForPromises()`

Blocks until all promises created have completed.

## License

MIT
