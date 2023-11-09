package gol

var GolHandler = "GolOperations.Update"
var BrokerHandler = "BrokerOperations.Execute"
var TickerHandler = "BrokerOperations.AliveCells"
var SaveHandler = "BrokerOperations.Save"

type BrokerResponse struct {
	World [][]byte
}

type DistributorRequest struct {
	P     Params
	World [][]byte
}

type BrokerRequest struct {
	P     Params
	World [][]byte
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
