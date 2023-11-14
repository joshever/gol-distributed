package gol

func Next(p Params, world [][]byte, update chan [][]byte, startY int, endY int) {
	// Sequential if 1 thread
	world = calculateNextState(p, world, startY, endY)
	update <- world
	return
}

func calculateNextState(p Params, world [][]byte, startY, endY int) [][]byte {
	// Deep Copy each world to avoid shared memory
	newWorld := MakeNewWorld(p, world)
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

func CreateEmptyWorld(p Params) [][]byte {
	world := make([][]byte, p.ImageHeight)
	for k := range world {
		world[k] = make([]byte, p.ImageWidth)
	}
	return world
}

func MakeNewWorld(p Params, world [][]byte) [][]byte {
	newWorld := CreateEmptyWorld(p)
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
