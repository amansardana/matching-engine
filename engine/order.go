package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/amansardana/matching-engine/types"
	"github.com/amansardana/matching-engine/utils"
	"github.com/gomodule/redigo/redis"
)

type FillOrder struct {
	Amount uint64
	Order  *types.Order
}

func (e *EngineResource) matchOrder(order *types.Order) (err error) {
	var match *Match
	if order.Type == types.SELL {
		match, err = e.sellOrder(order)
	} else if order.Type == types.BUY {
		match, err = e.buyOrder(order)
	}
	mab, err := json.Marshal(match)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("\n======>\n%s\n<======\n", mab)
	// Note: Plug the option for orders like FOC, Limit, OnlyFill (If Required)

	// If NO_MATCH add to order book
	if match.FillStatus == NO_MATCH {
		e.addOrder(order)
		return
	}

	// Execute Trade
	engineResponse, err := e.execute(match)
	if err != nil {
		log.Printf("\nexecute XXXXXXX\n%s\nXXXXXXX execute\n", err)
	}
	erab, err := json.Marshal(engineResponse)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("\n======> engineResponse\n%s\nengineResponse <======\n", erab)
	return
}

func (e *EngineResource) buyOrder(order *types.Order) (match *Match, err error) {
	match = &Match{
		Order:      order,
		FillStatus: NO_MATCH,
	}
	match.MatchingOrders = make([]*FillOrder, 0)

	oskv := order.GetOBMatchKey()

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
	fmt.Printf("\n======> priceRange\n%s\npriceRange <======\n", priceRange)

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

			for _, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return nil, err
				}

				match.FillStatus = PARTIAL

				// update filledAmount
				beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
				if beAmtAvailable > order.Amount {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{order.Amount, &bookEntry})
					filledAmount = order.Amount
				} else {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{beAmtAvailable, &bookEntry})
					filledAmount = beAmtAvailable
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
	match.MatchingOrders = make([]*FillOrder, 0)

	obkv := order.GetOBMatchKey()
	fmt.Println(obkv)
	fmt.Println("[" + utils.UintToPaddedString(order.Price))
	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZREVRANGEBYLEX", obkv, "+", "["+utils.UintToPaddedString(order.Price))) // "ZREVRANGEBYLEX" key max min
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

			for _, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return nil, err
				}

				match.FillStatus = PARTIAL

				// update filledAmount
				beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
				if beAmtAvailable > order.Amount {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{order.Amount, &bookEntry})
					filledAmount = order.Amount
				} else {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{beAmtAvailable, &bookEntry})
					filledAmount = beAmtAvailable
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
	res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
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
