package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
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

func (q *GolOperations) Shutdown(req gol.ShutdownRequest, res *gol.ShutdownResponse) (err error) {
	fmt.Println("Quitting Node...")
	os.Exit(0)
	return
}

func main() {
	rpc.Register(&GolOperations{})
	BrokerAddr := flag.String("broker", "8040", "Broker Listener Port")
	flag.Parse()
	listener, _ := net.Listen("tcp", ":"+*BrokerAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
