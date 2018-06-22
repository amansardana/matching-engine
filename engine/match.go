package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/amansardana/matching-engine/utils"

	"github.com/amansardana/matching-engine/types"
	"github.com/gomodule/redigo/redis"
)

type Match struct {
	Order          *types.Order
	FillStatus     FillStatus
	MatchingOrders []*FillOrder
}
type FillStatus int

type EngineResponse struct {
	Order          *types.Order
	Trades         []*types.Trade
	RemainingOrder *types.Order

	FillStatus     FillStatus
	MatchingOrders []*FillOrder
}

const (
	_ FillStatus = iota
	NO_MATCH
	PARTIAL
	FULL
)

func (e *EngineResource) execute(m *Match) (response *EngineResponse, err error) {

	if m == nil {
		err = errors.New("No match passed")
		return
	}

	trades := make([]*types.Trade, 0)
	var filledAmount uint64

	order := m.Order
	MatchedOrders := m.MatchingOrders
	remainigOrder := *order

	for i, o := range MatchedOrders {
		mo := o.Order
		ss, list := mo.GetOBKeys()
		// POP the order from the top of list
		reply, err := redis.Bytes(e.redisConn.Do("LPOP", list)) // "ZREVRANGEBYLEX" key max min
		if err != nil {
			log.Printf("LPOP: %s\n", err)
			return nil, err
		}

		var bookEntry types.Order
		err = json.Unmarshal(reply, &bookEntry)
		if err != nil {
			log.Printf("json.Unmarshal: %s\n", err)
			return nil, err
		}

		if bookEntry.ID != mo.ID {
			log.Fatal("Invalid matching order passed: ", bookEntry.ID, mo.ID, list)
			return nil, errors.New("Invalid matching order passed")
		}
		filledAmount = filledAmount + o.Amount
		remainigOrder.Amount = remainigOrder.Amount - o.Amount

		// Create trade object to be passed to the system for further processing
		t := &types.Trade{
			Amount:     o.Amount,
			Price:      order.Price,
			OrderHash:  mo.Hash,
			Type:       order.Type,
			TradeNonce: uint64(i),
			Taker:      order.UserAddress,
			PairName:   order.PairName,
		}
		// TODO: Implement compute hash functions
		// t.Hash = t.ComputeHash()

		trades = append(trades, t)

		// If book entry order is not filled completely then update the filledAmount and push it back to the head of list

		if (bookEntry.Amount - bookEntry.FilledAmount) > o.Amount {
			bookEntry.FilledAmount = bookEntry.FilledAmount + o.Amount
			bookEntryAsBytes, err := json.Marshal(bookEntry)
			if err != nil {
				log.Printf("json.Marshal: %s", err)
			}
			res, err := e.redisConn.Do("LPUSH", list, bookEntryAsBytes)
			if err != nil {
				log.Printf("LPUSH: %s", err)
			}
			fmt.Println("LPUSH: ", res)
		}

		// Get length of remaining orders in the list

		l, err := redis.Uint64(e.redisConn.Do("LLEN", list))
		if err != nil {
			log.Printf("LLEN: %s", err)
		} else if l == 0 {
			// If list is empty: remove the list and remove the price point from sorted set
			_, err := e.redisConn.Do("del", list)
			if err != nil {
				log.Printf("del: %s", err)
				return nil, err
			}
			// fmt.Printf("del: %s", res)

			_, err = e.redisConn.Do("ZREM", ss, utils.UintToPaddedString(mo.Price))
			if err != nil {
				log.Printf("ZREM: %s", err)
				return nil, err
			}
			// fmt.Printf("ZREM: %s", res)
		}
	}

	order.FilledAmount = filledAmount

	response = &EngineResponse{
		Order:          order,
		Trades:         trades,
		FillStatus:     m.FillStatus,
		MatchingOrders: m.MatchingOrders,
	}
	if remainigOrder.Amount != 0 {
		response.RemainingOrder = &remainigOrder
	}
	return
}
