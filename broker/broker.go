package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"uk.ac.bris.cs/gameoflife/gol"
)

var (
	alive     int
	turns     int
	mutex     sync.Mutex
	World     [][]byte
	quit      bool
	shutDown  bool
	nodeAddrs []string
	nodes     []*rpc.Client
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
	// Connect to nodes
	nodes = make([]*rpc.Client, len(nodeAddrs))
	for i, addr := range nodeAddrs {
		nodes[i], _ = rpc.Dial("tcp", "127.0.0.1:"+addr)
	}
	// Initialise world, p and strip size
	p := req.P
	world := req.World
	height := p.ImageHeight / len(nodes)
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
			for _, node := range nodes {
				node.Call(gol.ShutNodeHandler, new(gol.ShutdownRequest), new(gol.ShutdownResponse))
			}
			fmt.Println("Quitting Broker...")
			mutex.Unlock()
			os.Exit(0)
			return
		}

		// Make lists of requests, Responses and Done calls
		requests := make([]gol.BrokerRequest, len(nodes))
		responses := make([]gol.NodeResponse, len(nodes))
		done := make([]chan *rpc.Call, len(nodes))

		// Call nodes to calculate next
		fmt.Println("Executing turn", i+1)
		world = callNodes(p, nodes, requests, responses, done, world, height)

		// Update global data
		World = world
		alive = len(gol.CalculateAliveCells(World))
		turns = i + 1
		mutex.Unlock()
	}
	res.World = world
	for _, node := range nodes {
		node.Close()
	}
	return
}

func callNodes(p gol.Params, nodes []*rpc.Client, requests []gol.BrokerRequest, responses []gol.NodeResponse, done []chan *rpc.Call, world [][]byte, height int) [][]byte {
	for j, node := range nodes {
		responses[j] = *new(gol.NodeResponse)
		requests[j] = *new(gol.BrokerRequest)
		requests[j].P = p
		requests[j].World = world
		requests[j].StartY = j * height
		if j == len(nodes)-1 {
			requests[j].EndY = p.ImageHeight
		} else {
			requests[j].EndY = (j + 1) * height
		}
		done[j] = make(chan *rpc.Call, 1)
		node.Go(gol.GolHandler, requests[j], &responses[j], done[j])
	}
	var newWorld [][]byte
	for k := range nodes {
		<-done[k]
		startY := requests[k].StartY
		endY := requests[k].EndY
		newWorld = append(newWorld, responses[k].World[startY:endY]...)
	}
	return newWorld
}

func main() {
	// Parse node addresses from command line
	var stringList string
	flag.StringVar(&stringList, "stringList", "", "Comma-separated list of strings")
	flag.Parse()
	nodeAddrs = strings.Split(stringList, ",")
	fmt.Println("Node Addresses: " + stringList)
	// listen for distributor connection
	rpc.Register(&BrokerOperations{})
	listener, _ := net.Listen("tcp", ":"+"8030")
	defer listener.Close()
	rpc.Accept(listener)
}
