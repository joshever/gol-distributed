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
	go gol.Next(p, world, update, req.StartY, req.EndY)
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
	var port string
	flag.StringVar(&port, "port", "8070", "Node Listener Port")
	flag.Parse()
	fmt.Println("Node port:", port)
	listener, _ := net.Listen("tcp", ":"+port)
	defer listener.Close()
	rpc.Register(&GolOperations{})
	rpc.Accept(listener)
}
