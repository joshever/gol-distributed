package gol

import (
	"fmt"
	"net/rpc"
	"os"
)

// Key Presses Function
func keyPresses(p Params, c distributorChannels, broker *rpc.Client, keyPressesDone chan bool) {
	//pause := false
	for {
		select {
		case <-keyPressesDone:
			return
		case x := <-c.keys:
			switch x {
			case 's':
				request := new(SaveRequest)
				response := new(SaveResponse)
				saveError := broker.Call(SaveHandler, request, response)
				Handle(saveError)
				writePgm(p, c, response.World, response.Turns)
				clearKeys(c)
			case 'q':
				request := new(QuitRequest)
				response := new(QuitResponse)
				quitError := broker.Call(QuitHandler, request, response)
				Handle(quitError)
				c.ioCommand <- ioCheckIdle
				<-c.ioIdle
				c.events <- StateChange{response.Turns, Quitting}
				os.Exit(0)
			case 'k':
				request := new(ShutdownRequest)
				response := new(ShutdownResponse)
				fmt.Println("Quitting Distributor...")
				shutDownError := broker.Call(ShutBrokerHandler, request, response)
				Handle(shutDownError)
				writePgm(p, c, response.World, response.Turns)
				c.ioCommand <- ioCheckIdle
				<-c.ioIdle
				c.events <- StateChange{response.Turns, Quitting}
				os.Exit(0)
			}
		}
	}
}

func clearKeys(c distributorChannels) {
	for {
		select {
		case <-c.keys:
		default:
			return
		}
	}
}
