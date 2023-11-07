package gol

func next(p Params, world [][]byte, update chan [][]byte) {
	// Sequential if 1 thread
	if p.Threads == 1 {
		world = calculateNextState(p, world, 0, p.ImageHeight)
	} else {
		world = parallel(p, world)
	}
	update <- world
	return
}

func calculateNextState(p Params, world [][]byte, startY, endY int) [][]byte {
	// Deep Copy each world to avoid shared memory
	newWorld := makeNewWorld(p, world)
	for j := startY; j < endY; j++ {
		for i := 0; i < p.ImageWidth; i++ {
			// Use un-copied world as operations are read only
			aliveNeighbours := findAliveNeighbours(p, world, j, i)
			if world[j][i] == ALIVE {
				if aliveNeighbours < 2 {
					newWorld[j][i] = DEAD
				} else if aliveNeighbours <= 3 {
					newWorld[j][i] = ALIVE
				} else {
					newWorld[j][i] = DEAD
				}
			} else {
				if aliveNeighbours == 3 {
					newWorld[j][i] = ALIVE
				}
			}
		}
	}
	return newWorld
}

func createEmptyWorld(p Params) [][]byte {
	world := make([][]byte, p.ImageHeight)
	for k := range world {
		world[k] = make([]byte, p.ImageWidth)
	}
	return world
}

func makeNewWorld(p Params, world [][]byte) [][]byte {
	newWorld := createEmptyWorld(p)
	for k := range world {
		copy(newWorld[k], world[k])
	}
	return newWorld
}

func findAliveNeighbours(p Params, world [][]byte, x int, y int) int {
	alive := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			ny, nx := y+j, x+i
			if ny == p.ImageHeight {
				ny = 0
			} else if ny < 0 {
				ny = p.ImageHeight - 1
			}
			if nx < 0 {
				nx = p.ImageWidth - 1
			} else if nx == p.ImageWidth {
				nx = 0
			}
			if world[nx][ny] == ALIVE {
				alive += 1
			}
		}
	}
	return alive
}

func parallel(p Params, world [][]byte) [][]byte {
	var newPixelData [][]byte
	newHeight := p.ImageHeight / p.Threads
	// List of channels for each thread
	channels := make([]chan [][]byte, p.Threads)
	for i := 0; i < p.Threads; i++ {
		channels[i] = make(chan [][]byte)
		// Cover gaps missed due to rounding in last strip
		if i == p.Threads-1 {
			go worker(p, i*newHeight, p.ImageHeight, world, channels[i])
		} else {
			go worker(p, i*newHeight, (i+1)*newHeight, world, channels[i])
		}
	}
	for i := 0; i < p.Threads; i++ {
		// Read from specific channels in order to reassemble
		newPixelData = append(newPixelData, <-channels[i]...)
	}
	return newPixelData
}

func worker(p Params, startY, endY int, world [][]byte, out chan<- [][]uint8) {
	// Pass whole world which is then deep copied
	returned := calculateNextState(p, world, startY, endY)
	// Slice output into correct strip
	returned = returned[startY:endY]
	out <- returned
}
