package contract

// ExecutorFunc is the function passed to the promise, which performs the async operation.
// It receives resolve and reject functions to signal completion or failure.
type ExecutorFunc[T any] func(resolve func(T), reject func(error), finally func())
