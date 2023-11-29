package gol

var PauseHandler = "BrokerOperations.Pause"
var UnpauseHandler = "BrokerOperations.Unpause"
var QuitHandler = "BrokerOperations.Quit"
var ShutBrokerHandler = "BrokerOperations.Shutdown"
var SaveHandler = "BrokerOperations.Save"
var TickerHandler = "BrokerOperations.AliveCells"
var BrokerHandler = "BrokerOperations.Execute"
var GolHandler = "GolOperations.Update"
var ShutNodeHandler = "GolOperations.Shutdown"
var SDLHandler = "DistributorOperations.SDL"

type BrokerResponse struct {
	World [][]byte
}

type DistributorRequest struct {
	P     Params
	World [][]byte
}

type BrokerRequest struct {
	StartY int
	EndY   int
	P      Params
	World  [][]byte
}

type NodeResponse struct {
	World [][]byte
}

type AliveRequest struct {
}

type AliveResponse struct {
	Alive int
	Turns int
}

type SaveRequest struct{}

type SaveResponse struct {
	World [][]byte
	Turns int
}

type QuitRequest struct{}

type QuitResponse struct {
	Turns int
}

type PauseRequest struct{}

type PauseResponse struct {
	Turns int
}

type ShutdownRequest struct{}

type ShutdownResponse struct {
	World [][]byte
	Turns int
}

type SDLRequest struct {
	CellsFlipped [][]byte
	Turns        int
}
type SDLResponse struct {
}
