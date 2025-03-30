package counter

import (
	"log"
	"sync"
)

type ObservableCounter struct {
	sync.RWMutex
	counter   int
	observers map[chan int]struct{}
}

func NewObservableCounter() *ObservableCounter {
	return &ObservableCounter{
		observers: make(map[chan int]struct{}),
	}
}

func (oc *ObservableCounter) Subscribe() (chan int, int) {
	ch := make(chan int, 10)

	oc.Lock()
	oc.observers[ch] = struct{}{}
	currentCount := oc.counter
	oc.Unlock()

	return ch, currentCount
}

func (oc *ObservableCounter) Unsubscribe(ch chan int) {
	oc.Lock()
	delete(oc.observers, ch)
	close(ch)
	oc.Unlock()
}

func (oc *ObservableCounter) Inc() {
	oc.Lock()
	oc.counter++
	newVal := oc.counter
	oc.Unlock()

	oc.RLock()
	for ch := range oc.observers {
		select {
		case ch <- newVal:
		default:
			log.Printf("Observer channel buffer full for %v, skipping update", ch)
		}
	}
	oc.RUnlock()
}

func (oc *ObservableCounter) GetValue() int {
	oc.RLock()
	defer oc.RUnlock()

	return oc.counter
}
