package gol

import (
	"fmt"
	"os"
	"sync"
)

// Key Presses Function
func keyPresses(p Params, w *World, c distributorChannels, keyPressesDone chan bool, pauseDistributor chan bool, pauseTicker chan bool, mutex *sync.Mutex) {
	pause := false
	for {
		select {
		case <-keyPressesDone:
			return
		case x := <-c.keys:
			switch x {
			case 's':
				mutex.Lock()
				writePgm(p, c, w)
				clearKeys(c)
				mutex.Unlock()
			case 'q':
				mutex.Lock()
				writePgm(p, c, w)
				c.ioCommand <- ioCheckIdle
				<-c.ioIdle
				c.events <- StateChange{w.turns, Quitting}
				os.Exit(0)
			case 'p':
				pauseDistributor <- true
				pauseTicker <- true
				if pause == false {
					mutex.Lock()
					fmt.Println(fmt.Sprintf("Currently processing: %d", w.turns))
					mutex.Unlock()
					pause = true
				} else {
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
