package main

import (
	"fmt"
	"time"

	"promise/pkg"
)

// main function demonstrates the usage of the Promise implementation.
func main() {
	fmt.Println("--- Running Success Case ---")

	// --- Success Example ---
	// Create a new promise that resolves with a string.
	pkg.NewPromise[string](func(resolve func(string), reject func(error), finally func()) {
		// Simulate an async operation like an API call
		fmt.Println("Success Promise: Starting async work...")
		time.Sleep(2 * time.Second)
		resolve("Promise resolved successfully! Data has arrived.")
	}).Then(func(result string) {
		// This block runs on success
		fmt.Printf("Success Promise Then: %s\n", result)
	}).Then(func(result string) {
		// This block runs on success
		fmt.Printf("Success Again Promise Then: %s\n", result)
	}).Catch(func(err error) {
		// This block runs on failure
		fmt.Printf("Success Promise Catch: %s\n", err.Error())
	})

	fmt.Println("\n--- Running Failure Case ---")
	// --- Failure Example ---
	// Create a new promise that is expected to fail.
	promise := pkg.NewPromise[int](func(resolve func(int), reject func(error), finally func()) {
		// Simulate a failing async operation
		fmt.Println("Failure Promise: Starting async work...")
		time.Sleep(2 * time.Second)
		reject(fmt.Errorf("something went wrong during the integer operation"))
	})
	promise.Then(func(result int) {
		// This block won't run in this case
		fmt.Printf("Failure Promise Then: Received %d\n", result)
	})
	promise.Catch(func(err error) {
		// This block will run
		fmt.Printf("Failure Promise Catch: %s\n", err.Error())
	})
	promise.Finally(func() {
		// This block runs regardless of success or failure
		fmt.Println("Failure Promise Finally: Cleanup after failure.")
	})

	// Wait for the promises to complete
	fmt.Println("\nMain function is waiting for promises to complete...")
	pkg.WaitForPromises()

	fmt.Println("\nAll promises have completed. Exiting.")
}
