package gol

import (
	"fmt"
	"net/rpc"
	"os"
)

// Key Presses Function
func keyPresses(p Params, c distributorChannels, broker *rpc.Client, keyPressesDone chan bool, pauseTicker chan bool, tickerDone chan bool) {
	pause := false
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
				tickerDone <- true
				response := new(SaveResponse)
				fmt.Println("Quitting Distributor...")
				saveError := broker.Call(SaveHandler, new(SaveRequest), response)
				Handle(saveError)
				writePgm(p, c, response.World, response.Turns)
				shutDownError := broker.Call(ShutBrokerHandler, new(ShutdownRequest), new(ShutdownResponse))
				Handle(shutDownError)
				c.ioCommand <- ioCheckIdle
				<-c.ioIdle
				c.events <- StateChange{response.Turns, Quitting}
				os.Exit(0)
			case 'p':
				pauseTicker <- true
				response := new(PauseResponse)
				if pause == false {
					broker.Call(PauseHandler, new(PauseRequest), response)
					fmt.Println("Currently Executing", response.Turns)
					pause = true
				} else {
					broker.Call(UnpauseHandler, new(PauseRequest), response)
					fmt.Println("Continuing")
					pause = false
				}
				clearKeys(c)
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
