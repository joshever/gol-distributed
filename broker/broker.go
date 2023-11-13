package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"uk.ac.bris.cs/gameoflife/gol"
)

var (
	alive    int
	turns    int
	mutex    sync.Mutex
	World    [][]byte
	quit     bool
	shutDown bool
)

type BrokerOperations struct{}

func (b *BrokerOperations) Pause(req gol.PauseRequest, res *gol.PauseResponse) (err error) {
	mutex.Lock()
	res.Turns = turns
	return
}

func (b *BrokerOperations) Unpause(req gol.PauseRequest, res *gol.PauseResponse) (err error) {
	mutex.Unlock()
	return
}

func (b *BrokerOperations) Quit(req gol.QuitRequest, res *gol.QuitResponse) (err error) {
	mutex.Lock()
	res.Turns = turns
	quit = true
	mutex.Unlock()
	return
}

func (b *BrokerOperations) Shutdown(req gol.ShutdownRequest, res *gol.ShutdownResponse) (err error) {
	mutex.Lock()
	shutDown = true
	mutex.Unlock()
	return
}

func (b *BrokerOperations) Save(req gol.SaveRequest, res *gol.SaveResponse) (err error) {
	mutex.Lock()
	res.World = World
	res.Turns = turns
	mutex.Unlock()
	return
}

func (b *BrokerOperations) AliveCells(req gol.AliveRequest, res *gol.AliveResponse) (err error) {
	mutex.Lock()
	res.Alive = alive
	res.Turns = turns
	mutex.Unlock()
	return
}

func (b *BrokerOperations) Execute(req gol.DistributorRequest, res *gol.BrokerResponse) (err error) {
	// Connect to node
	node, dialErr := rpc.Dial("tcp", "127.0.0.1:8040")
	gol.Handle(dialErr)
	defer node.Close()

	// Initialise world, p and request/response
	p := req.P
	world := req.World
	response := new(gol.NodeResponse)
	request := new(gol.BrokerRequest)
	request.P = p

	// Initialise globals
	mutex.Lock()
	World = world
	turns = 0
	alive = len(gol.CalculateAliveCells(World))
	mutex.Unlock()

	// Return world to distributor if no turns
	if p.Turns == 0 {
		res.World = req.World
		return
	}

	// Call node to carry out each turn and return when done
	for i := 0; i < p.Turns; i++ {
		// Checks global quit, shutdown, pause
		mutex.Lock()
		if quit {
			quit = false
			fmt.Println("Resetting state..")
			res.World = world
			mutex.Unlock()
			return
		} else if shutDown {
			shutDown = false
			node.Call(gol.ShutNodeHandler, new(gol.ShutdownRequest), new(gol.ShutdownResponse))
			fmt.Println("Quitting Broker...")
			mutex.Unlock()
			os.Exit(0)
			return
		}

		// Call node to calculate next
		fmt.Println("Executing turn", i+1)
		request.World = world
		nodeErr := node.Call(gol.GolHandler, request, response)
		gol.Handle(nodeErr)
		world = response.World

		// Update global data
		World = world
		alive = len(gol.CalculateAliveCells(World))
		turns = i + 1
		mutex.Unlock()
	}
	res.World = response.World
	return
}

func main() {
	rpc.Register(&BrokerOperations{})
	listener, _ := net.Listen("tcp", ":"+"8030")
	defer listener.Close()
	rpc.Accept(listener)
}
