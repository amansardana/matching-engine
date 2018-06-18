package engine

import (
	"fmt"

	"github.com/amansardana/matching-engine/types"
)

func (e *EngineResource) buyOrder(order *types.Order) (err error) {
	okv := order.GetKVPrefix() + "/buy"
	fmt.Println(okv)
	// e.redisConn.Do("ZINCRBY", okv, order.Price, order.Amount)
	// TODO: ReThink the possibility of using sorted sets
	rply, err := e.redisConn.Do("ZRANGEBYSCORE", okv, order.Price-1, order.Price+1)

	// GET Range of sellOrder between minimum Sell order and order.Price

	// If no entry
	// ZINCRBY (Add order to order book)

	// else
	// check score(volume) at particular price
	// Can be filled completely?
	// Get orders from order list
	// make trade and add remaining orders to orderlist

	fmt.Printf("%s", rply)
	return
}

func (e *EngineResource) sellOrder(order *types.Order) (err error) {
	okv := order.GetKVPrefix() + "/sell"
	e.redisConn.Do("ZINCRBY", okv, order.Price, order.Amount)
	fmt.Println(e.redisConn.Do("ZRANGEBYSCORE", okv, order.Amount-1, order.Amount+1))
	return
}
