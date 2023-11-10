package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
	"uk.ac.bris.cs/gameoflife/gol"
)

var alive int
var turns int
var mutex sync.Mutex
var World [][]byte
var quit bool
var shutDown bool

type BrokerOperations struct{}

func (b *BrokerOperations) Quit(req gol.QuitRequest, res *gol.QuitResponse) (err error) {
	mutex.Lock()
	res.Turns = turns
	quit = true
	mutex.Unlock()
	return
}

func (b *BrokerOperations) Shutdown(req gol.ShutdownRequest, res *gol.ShutdownResponse) (err error) {
	mutex.Lock()
	res.Turns = turns
	res.World = World
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
	node, dialErr := rpc.Dial("tcp", "127.0.0.1:8040")
	gol.Handle(dialErr)
	defer node.Close()

	// Initialise world, p and request/response
	response := new(gol.NodeResponse)
	request := new(gol.BrokerRequest)
	p := req.P
	world := req.World
	request.P = p

	// Return world to distributor if no turns
	if p.Turns == 0 {
		res.World = world
		return
	}

	// Call node to carry out each turn and return when done
	for i := 0; i < p.Turns; i++ {
		if quit {
			quit = false
			fmt.Println("Resetting state..")
			return
		}
		if shutDown {
			node.Call(gol.ShutNodeHandler, new(gol.ShutdownRequest), new(gol.ShutdownResponse))
			fmt.Println("Quitting Broker...")
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}
		fmt.Println("Executing turn", i)
		// Update interface data
		mutex.Lock()
		World = world
		alive = len(gol.CalculateAliveCells(world))
		turns = i
		mutex.Unlock()
		// Call node to calculate next
		request.World = world
		nodeErr := node.Call(gol.GolHandler, request, response)
		gol.Handle(nodeErr)
		world = response.World
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
