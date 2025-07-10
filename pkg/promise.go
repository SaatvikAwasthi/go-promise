package pkg

import (
	"sync"

	"promise/pkg/contract"
)

// Promise is a generic struct that represents the eventual completion (or failure)
// of an asynchronous operation and its resulting value.
// It uses a sync.Mutex to handle concurrent access to its handlers,
// making it safe for cases where .Then or .Catch might be called after resolution.
type Promise[T any] struct {
	mutex   sync.Mutex
	then    func(T)
	catch   func(error)
	finally func()
}

// NewPromise creates and returns a new Promise.
// It takes an executor function that will be run in a separate goroutine.
func NewPromise[T any](executor contract.ExecutorFunc[T]) *Promise[T] {
	p := &Promise[T]{}
	wg.Add(1)

	// The resolve function handles the successful completion of the promise.
	resolve := func(value T) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if p.then != nil {
			// We launch the handler in a new goroutine to avoid blocking the
			// original executor goroutine if the .Then handler is slow.
			go func() {
				p.then(value)
				wg.Done()
			}()
		}
	}

	// The reject function handles the failure of the promise.
	reject := func(err error) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if p.catch != nil {
			// Same as resolve, run handler in a new goroutine.
			go func() {
				p.catch(err)
				wg.Done()
			}()
		}
	}

	finally := func() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if p.finally != nil {
			go func() {
				p.finally()
				wg.Done()
			}()
		}
	}

	// The core of the async operation. We run the executor in a new goroutine
	// so that the NewPromise call doesn't block.
	go executor(resolve, reject, finally)

	return p
}

// Then sets the success handler for the promise.
// It returns the promise itself to allow for chaining `Catch`.
func (p *Promise[T]) Then(handler func(T)) *Promise[T] {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.then = handler
	return p
}

// Catch sets the error handler for the promise.
// It returns the promise itself to allow for chaining `Finally`.
func (p *Promise[T]) Catch(handler func(error)) *Promise[T] {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.catch = handler
	return p
}

// Finally adds a handler that will be called regardless of whether the promise
// resolves or rejects.
func (p *Promise[T]) Finally(handler func()) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	oldThen := p.then
	oldCatch := p.catch

	if oldThen != nil {
		p.then = func(value T) {
			oldThen(value)
			handler()
		}
	} else {
		p.then = func(_ T) {
			handler()
		}
	}

	if oldCatch != nil {
		p.catch = func(err error) {
			oldCatch(err)
			handler()
		}
	} else {
		p.catch = func(_ error) {
			handler()
		}
	}
}
