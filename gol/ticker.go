package gol

import (
	"sync"
	"time"
)

// Ticker function
func ticker(w *World, c distributorChannels, tickerDone chan bool, pauseTicker chan bool, mutex *sync.Mutex) {
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-pauseTicker:
			<-pauseTicker
		case <-tickerDone:
			return
		case <-tick.C:
			// Mutex to cover data race for w
			mutex.Lock()
			c.events <- AliveCellsCount{w.turns, len(calculateAliveCells(w.world))}
			mutex.Unlock()
		}
	}
}
