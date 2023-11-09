package gol

import (
	"net/rpc"
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
				//case 'q':
				//	writePgm(p, c, w)
				//	c.ioCommand <- ioCheckIdle
				//	<-c.ioIdle
				//	c.events <- StateChange{w.turns, Quitting}
				//	os.Exit(0)
				//case 'p':
				//	if pause == false {
				//		fmt.Println(fmt.Sprintf("Currently processing: %d", w.turns))
				//		pause = true
				//	} else {
				//		fmt.Println("Continuing")
				//		pause = false
				//	}
				//	clearKeys(c)
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
