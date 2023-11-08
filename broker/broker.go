package main

import (
	"flag"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/gol"
)

type BrokerOperations struct{}

func (b *BrokerOperations) AliveCells(req gol.AliveRequest, res *gol.AliveResponse) (err error) {
	return
}

func (b *BrokerOperations) Execute(req gol.DistributorRequest, res *gol.BrokerResponse) (err error) {
	flag.Parse()
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
	DistAddr := flag.String("DistAddr", "8030", "Dist Listener Port")
	flag.Parse()
	listener, _ := net.Listen("tcp", ":"+*DistAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
