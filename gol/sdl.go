package gol

import (
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

func sdl(p Params, c distributorChannels, broker *rpc.Client, sdlDone chan bool, old [][]byte) {
	tick := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-sdlDone:
			return
		case <-tick.C:
			request := new(SaveRequest)
			response := new(SaveResponse)
			saveError := broker.Call(SaveHandler, request, response)
			Handle(saveError)
			world := response.World
			turns := response.Turns
			for j := 0; j < p.ImageHeight; j++ {
				for i := 0; i < p.ImageWidth; i++ {
					if world[j][i] != old[j][i] {
						c.events <- CellFlipped{turns, util.Cell{i, j}}
					}
				}
			}
			c.events <- TurnComplete{turns}
			old = world
		}
	}
}
