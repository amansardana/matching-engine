package services

import (
	"errors"
	"fmt"

	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/engine"
	"github.com/amansardana/matching-engine/types"
)

type OrderService struct {
	orderDao   *daos.OrderDao
	balanceDao *daos.BalanceDao
	pairDao    *daos.PairDao
}

func NewOrderService(orderDao *daos.OrderDao, balanceDao *daos.BalanceDao, pairDao *daos.PairDao) *OrderService {
	return &OrderService{orderDao, balanceDao, pairDao}
}

func (s *OrderService) Create(order *types.Order) (err error) {

	// Fill token and pair data

	p, err := s.pairDao.GetByName(order.PairName)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("Pair not found")
	}
	order.PairID = p.ID
	order.BuyToken = p.BuyTokenSymbol
	order.BuyTokenAddress = p.BuyTokenSymbol
	order.SellToken = p.SellTokenSymbol
	order.SellTokenAddress = p.SellTokenAddress

	// TODO: Validate if order is valid

	// Validate user has balance
	bal, err := s.balanceDao.GetByAddress(order.UserAddress)
	if err != nil {
		return err
	}
	if order.Type == types.BUY {
		amt := bal.Tokens[order.SellToken]
		if amt.Amount < order.AmountSell+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Buy : Verified")

		amt.Amount = amt.Amount - (order.AmountSell + order.Fee)
		amt.LockedAmount = amt.Amount + (order.AmountSell + order.Fee)
		err = s.balanceDao.LockFunds(order.UserAddress, order.SellToken, &amt)

		if err != nil {
			return err
		}

	} else if order.Type == types.SELL {
		amt := bal.Tokens[order.BuyToken]
		if amt.Amount < order.AmountBuy+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Sell : Verified")
		amt.Amount = amt.Amount - (order.AmountBuy + order.Fee)
		amt.LockedAmount = amt.Amount + (order.AmountBuy + order.Fee)
		err = s.balanceDao.LockFunds(order.UserAddress, order.BuyToken, &amt)

		if err != nil {
			return err
		}
	}
	order.ID = bson.NewObjectId()

	// if err = s.orderDao.Create(order); err != nil {
	// 	return
	// }
	// TODO: Push order to queue
	engine.Engine.PublishOrder(order)
	return err
}

func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}
func (s *OrderService) GetAll() ([]types.Order, error) {
	return s.orderDao.GetAll()
}
