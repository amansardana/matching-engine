package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/amansardana/matching-engine/types"
	"github.com/amansardana/matching-engine/utils"
	"github.com/gomodule/redigo/redis"
)


func (e *EngineResource) matchOrder(order *types.Order) (err error) {
	var match *Match
	if order.Type == types.SELL {
		match, err = e.sellOrder(order)
	} else if order.Type == types.BUY {
		e.buyOrder(order)
	}
	mab, err := json.Marshal(match)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("\n======>\n%s\n<======\n", mab)
	// Note: Plug the option for orders like FOC, Limit, OnlyFill (If Required)

	// Execute Trade
	if match.FillStatus == NO_MATCH {
		e.addOrder(order)
		return
	}
	return
}

func (e *EngineResource) buyOrder(order *types.Order) (match *Match, err error) {
	match = &Match{
		Order:      order,
		FillStatus: NO_MATCH,
	}
	match.MatchingOrders = make(map[int]*types.Order)

	// okv := order.GetKVPrefix() + "::buy"
	oskv := order.GetKVPrefix() + "::sell"

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZRANGEBYLEX", oskv, "-", "["+utils.UintToPaddedString(order.Price))) // "ZRANGEBYLEX" key min max
	if err != nil {
		log.Printf("ZREVRANGEBYLEX: %s\n", err)
		return
	}

	priceRange := make([]uint64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
	}

	var filledAmount uint64
	var orderAmount = order.Amount

	if len(priceRange) == 0 {
		match.FillStatus = NO_MATCH
	} else {
		for _, pr := range priceRange {
			reply, err := redis.ByteSlices(e.redisConn.Do("LRANGE", oskv+"::"+utils.UintToPaddedString(pr), 0, -1)) // "ZREVRANGEBYLEX" key max min
			if err != nil {
				log.Printf("LRANGE: %s\n", err)
				return nil, err
			}

			for index, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return nil, err
				}

				match.MatchingOrders[index] = order
				match.FillStatus = PARTIAL

				// update filledAmount
				if bookEntry.Amount > order.Amount {
					filledAmount = order.Amount
				} else if bookEntry.Amount == order.Amount {
					filledAmount = bookEntry.Amount
				} else {
					filledAmount = bookEntry.Amount
				}

				if filledAmount == orderAmount {
					match.FillStatus = FULL
					// order filled return
					return match, nil
				}
			}
		}
	}
	return
}

func (e *EngineResource) sellOrder(order *types.Order) (match *Match, err error) {
	match = &Match{
		Order:      order,
		FillStatus: NO_MATCH,
	}
	match.MatchingOrders = make(map[int]*types.Order)

	// okv := order.GetKVPrefix() + "::sell"
	obkv := order.GetKVPrefix() + "::buy"

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZREVRANGEBYLEX", obkv, "["+utils.UintToPaddedString(order.Price), "-")) // "ZREVRANGEBYLEX" key max min
	if err != nil {
		log.Printf("ZREVRANGEBYLEX: %s\n", err)
		return
	}

	priceRange := make([]uint64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
	}

	var filledAmount uint64
	var orderAmount = order.Amount

	if len(priceRange) == 0 {
		match.FillStatus = NO_MATCH
	} else {
		for _, pr := range priceRange {
			reply, err := redis.ByteSlices(e.redisConn.Do("LRANGE", obkv+"::"+utils.UintToPaddedString(pr), 0, -1)) // "ZREVRANGEBYLEX" key max min
			if err != nil {
				log.Printf("LRANGE: %s\n", err)
				return nil, err
			}

			for index, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return nil, err
				}

				match.MatchingOrders[index] = order
				match.FillStatus = PARTIAL

				// update filledAmount
				if bookEntry.Amount > order.Amount {
					filledAmount = order.Amount
				} else if bookEntry.Amount == order.Amount {
					filledAmount = bookEntry.Amount
				} else {
					filledAmount = bookEntry.Amount
				}

				if filledAmount == orderAmount {
					match.FillStatus = FULL
					// order filled return
					return match, nil
				}
			}
		}
	}
	return
}

func (e *EngineResource) addOrder(order *types.Order) {

	ssKey, listKey := order.GetOBKeys()
	res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, order.Price) // Add price point to order book
	if err != nil {
		log.Printf("ZADD: %s", err)
	}
	fmt.Printf("ZADD: %s\n", res)

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("ZADD: %s", err)
	}
	res, err = e.redisConn.Do("RPUSH", listKey, orderAsBytes)
	if err != nil {
		log.Printf("RPUSH: %s", err)
	}
	fmt.Printf("RPUSH: %s\n", res)
}
