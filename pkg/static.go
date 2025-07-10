package pkg

import (
	"fmt"
	"sync"
)

// All waits for all promises to be resolved, or for any to be rejected.
// Returns a new Promise that resolves with an array of all results or rejects with the first error.
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	return NewPromise[[]T](func(resolve func([]T), reject func(error), finally func()) {
		if len(promises) == 0 {
			resolve([]T{})
			return
		}

		results := make([]T, len(promises))
		var mu sync.Mutex
		remaining := len(promises)

		for i, p := range promises {
			idx := i // Capture loop variable
			p.Then(func(val T) {
				mu.Lock()
				results[idx] = val
				remaining--
				if remaining == 0 {
					mu.Unlock()
					resolve(results)
				} else {
					mu.Unlock()
				}
			}).Catch(func(err error) {
				mu.Lock()
				// Only reject once
				if remaining > 0 {
					remaining = 0
					mu.Unlock()
					reject(err)
				} else {
					mu.Unlock()
				}
			})
		}
	})
}

// Race returns a promise that fulfills or rejects as soon as one of the promises fulfills or rejects.
func Race[T any](promises ...*Promise[T]) *Promise[T] {
	return NewPromise[T](func(resolve func(T), reject func(error), finally func()) {
		if len(promises) == 0 {
			reject(fmt.Errorf("no promises to race"))
			return
		}

		var settled bool
		var mu sync.Mutex

		for _, p := range promises {
			p.Then(func(val T) {
				mu.Lock()
				if !settled {
					settled = true
					mu.Unlock()
					resolve(val)
				} else {
					mu.Unlock()
				}
			}).Catch(func(err error) {
				mu.Lock()
				if !settled {
					settled = true
					mu.Unlock()
					reject(err)
				} else {
					mu.Unlock()
				}
			})
		}
	})
}

// PromiseResult represents the result of a promise that may be fulfilled or rejected
type PromiseResult[T any] struct {
	Value     T
	Error     error
	Fulfilled bool
}

// AllSettled waits until all promises have settled (either resolved or rejected).
// Returns a promise that resolves with an array of objects representing the settlement status of each promise.
func AllSettled[T any](promises ...*Promise[T]) *Promise[[]PromiseResult[T]] {
	return NewPromise[[]PromiseResult[T]](func(resolve func([]PromiseResult[T]), reject func(error), finally func()) {
		if len(promises) == 0 {
			resolve([]PromiseResult[T]{})
			return
		}

		results := make([]PromiseResult[T], len(promises))
		var mu sync.Mutex
		remaining := len(promises)

		for i, p := range promises {
			idx := i // Capture loop variable
			p.Then(func(val T) {
				mu.Lock()
				results[idx] = PromiseResult[T]{Value: val, Fulfilled: true}
				remaining--
				if remaining == 0 {
					mu.Unlock()
					resolve(results)
				} else {
					mu.Unlock()
				}
			}).Catch(func(err error) {
				mu.Lock()
				results[idx] = PromiseResult[T]{Error: err, Fulfilled: false}
				remaining--
				if remaining == 0 {
					mu.Unlock()
					resolve(results)
				} else {
					mu.Unlock()
				}
			})
		}
	})
}

// Any returns a promise that fulfills when any of the input promises fulfills, with this first fulfillment value.
// Rejects only if all promises reject, with an AggregateError containing all rejection reasons.
func Any[T any](promises ...*Promise[T]) *Promise[T] {
	return NewPromise[T](func(resolve func(T), reject func(error), finally func()) {
		if len(promises) == 0 {
			reject(fmt.Errorf("all promises rejected"))
			return
		}

		var mu sync.Mutex
		remaining := len(promises)
		errors := make([]error, len(promises))

		for i, p := range promises {
			idx := i // Capture loop variable
			p.Then(func(val T) {
				mu.Lock()
				if remaining > 0 {
					remaining = 0
					mu.Unlock()
					resolve(val)
				} else {
					mu.Unlock()
				}
			}).Catch(func(err error) {
				mu.Lock()
				errors[idx] = err
				remaining--
				if remaining == 0 {
					mu.Unlock()
					reject(fmt.Errorf("all promises rejected: %v", errors))
				} else {
					mu.Unlock()
				}
			})
		}
	})
}
