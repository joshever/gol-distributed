package main

import (
	"flag"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/gol"
)

type GolOperations struct{}

func (q *GolOperations) Update(req gol.BrokerRequest, res *gol.NodeResponse) (err error) {
	world := req.World
	p := req.P
	update := make(chan [][]byte)
	go gol.Next(p, world, update)
	world = <-update
	res.World = world
	return
}

func main() {
	BrokerAddr := flag.String("broker", "8040", "Broker Listener Port")
	flag.Parse()
	rpc.Register(&GolOperations{})
	listener, _ := net.Listen("tcp", ":"+*BrokerAddr)
	defer listener.Close()
	rpc.Accept(listener)
}