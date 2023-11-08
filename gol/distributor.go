package gol

import (
	"fmt"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/util"
)

const ALIVE = byte(255)
const DEAD = byte(0)

type World struct {
	world [][]byte
	turns int
}

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
	keys       <-chan rune
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	broker, dialErr := rpc.Dial("tcp", "127.0.0.1:8030")
	defer broker.Close()
	Handle(dialErr)
	world := setup(p, c)
	request := DistributorRequest{P: p, World: world}
	response := new(BrokerResponse)
	tickerDone := make(chan bool)
	go ticker(c, broker, tickerDone)
	brokerErr := broker.Call(BrokerHandler, request, response)
	Handle(brokerErr)
	world = response.World

	// Final Turn Complete
	writePgm(p, c, world)
	finalState := FinalTurnComplete{p.Turns, CalculateAliveCells(world)}
	c.events <- finalState
	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- StateChange{p.Turns, Quitting}
	tickerDone <- true
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

// Declare functions used in goroutines
func writePgm(p Params, c distributorChannels, world [][]byte) {
	outputFilename := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, p.Turns)
	c.ioCommand <- ioOutput
	c.ioFilename <- outputFilename
	for j := 0; j < p.ImageHeight; j++ {
		for i := 0; i < p.ImageWidth; i++ {
			c.ioOutput <- world[j][i]
		}
	}
}

func CalculateAliveCells(world [][]byte) []util.Cell {
	var cells = []util.Cell{}
	for j := range world {
		for i := range world[0] {
			if world[j][i] == byte(255) {
				cells = append(cells, util.Cell{i, j})
			}
		}
	}
	return cells
}

func setup(p Params, c distributorChannels) [][]byte {
	// Construct file name and trigger IO to fill channel with file bytes
	inputFilename := fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- inputFilename

	// Local turn and world variables
	// world is filled byte by byte from IO input
	world := createEmptyWorld(p)
	for j := 0; j < p.ImageHeight; j++ {
		for i := 0; i < p.ImageWidth; i++ {
			nextByte := <-c.ioInput
			world[j][i] = nextByte
			if nextByte == ALIVE {
				c.events <- CellFlipped{0, util.Cell{i, j}}
			}
		}
	}
	return world
}