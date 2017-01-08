package balancers

import (
	"sync"
	"testing"
)

func TestChoose(t *testing.T) {
	backends := []string{"a", "b", "c", "d"}
	counterMap := make(map[string]int)

	wg := new(sync.WaitGroup)
	lock := new(sync.RWMutex)
	iterations := 400
	wg.Add(iterations)

	for _, backend := range backends {
		counterMap[backend] = 0
	}
	rr := NewRoundRobin(backends)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			counterMap[rr.Choose()]++
		}()
	}
	wg.Wait()
	cond := counterMap["a"] == counterMap["b"] &&
		counterMap["b"] == counterMap["c"] &&
		counterMap["c"] == counterMap["d"] &&
		counterMap["d"] == iterations/len(backends)
	if !cond {
		t.Fatal("Count of hits to every backend should be equal to each other and to iterations/len(backends)!")
	}

}
