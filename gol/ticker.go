package gol

import (
	"net/rpc"
	"time"
)

// Ticker function
func ticker(c distributorChannels, broker *rpc.Client, tickerDone chan bool) {
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-tickerDone:
			return
		case <-tick.C:
			response := new(AliveResponse)
			brokerErr := broker.Call(TickerHandler, new(AliveRequest), response)
			Handle(brokerErr)
			c.events <- AliveCellsCount{response.Turns, response.Alive}
		}
	}
}
