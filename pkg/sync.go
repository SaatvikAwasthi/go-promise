package pkg

import (
	"sync"
)

// Wait for all promise to complete
var wg sync.WaitGroup

func WaitForPromises() {
	wg.Wait()
}
