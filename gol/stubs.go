package gol

var GolHandler = "GolOperations.Update"
var BrokerHandler = "BrokerOperations.Execute"

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
