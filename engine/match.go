package engine

import "github.com/amansardana/matching-engine/types"

type Match struct {
	Order          *types.Order
	FillStatus     FillStatus
	MatchingOrders map[int]*types.Order
}
type FillStatus int

const (
	_ FillStatus = iota
	NO_MATCH
	PARTIAL
	FULL
)

func (m *Match) CreateTrades() {
	
}