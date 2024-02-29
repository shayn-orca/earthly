package pcounter

import "sync"

var (
	mu             sync.Mutex
	activeRefFuncs int
	allowedRefFunc bool
)

func init() {
	allowedRefFunc = true
}

func CanDoRefFunc() (ok bool, release func()) {
	mu.Lock()
	defer mu.Unlock()
	if !allowedRefFunc {
		panic("ref funcs not allowed")
		//return false, func() {}
	}
	activeRefFuncs += 1
	releaserCalled := false
	return true, func() {
		mu.Lock()
		defer mu.Unlock()
		if releaserCalled {
			panic("releaser already called")
		}
		releaserCalled = true
		if !allowedRefFunc {
			panic("CanDoRefFunc released when allowedRefFunc=false") // shouldn't happen
		}
		activeRefFuncs -= 1
	}
}

func WaitUntilShutdownIsSafe() {
	mu.Lock()
	defer mu.Unlock()
	if activeRefFuncs > 0 {
		panic("shutting down while ref funcs are in progress")
	}
	allowedRefFunc = false
}
