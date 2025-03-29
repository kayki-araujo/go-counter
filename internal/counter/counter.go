package counter

import (
	"sync"
)

type ObservableCount struct {
	sync.RWMutex
	counter   int
	observers map[chan int]bool
}

func NewObservableCount() *ObservableCount {
	return &ObservableCount{
		observers: make(map[chan int]bool),
	}
}

func (oc *ObservableCount) Subscribe() (chan int, int) {
	ch := make(chan int, 10)

	oc.Lock()
	oc.observers[ch] = true
	currentCount := oc.counter
	oc.Unlock()

	return ch, currentCount
}

func (oc *ObservableCount) Unsubscribe(ch chan int) {
	oc.Lock()
	close(ch)
	delete(oc.observers, ch)
	oc.Unlock()
}

func (oc *ObservableCount) Inc() {
	oc.Lock()
	oc.counter++
	newVal := oc.counter
	oc.Unlock()

	oc.RLock()
	for ch := range oc.observers {
		select {
		case ch <- newVal:
		default:
		}
	}
	oc.RUnlock()
}

func (oc *ObservableCount) GetValue() int {
	oc.RLock()
	defer oc.RUnlock()
	return oc.counter
}
