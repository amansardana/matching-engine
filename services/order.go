package services

import (
	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/engine"
	"github.com/amansardana/matching-engine/types"
)

type OrderService struct {
	orderDao   *daos.OrderDao
	balanceDao *daos.BalanceDao
}

func NewOrderService(orderDao *daos.OrderDao, balanceDao *daos.BalanceDao) *OrderService {
	return &OrderService{orderDao, balanceDao}
}

func (s *OrderService) Create(order *types.Order) (err error) {

	// TODO: Validate if order is valid

	// Validate user has balance
	// bal, err := s.balanceDao.GetByAddress(order.UserAddress)
	// if err != nil {
	// 	return err
	// }
	// if order.Type == types.BUY {
	// 	amt := bal.Tokens[order.TokenSell]
	// 	if amt.Amount < order.AmountSell+order.Fee {
	// 		return errors.New("Insufficient Balance")
	// 	}
	// 	amt.Amount = amt.Amount - (order.AmountSell + order.Fee)
	// 	amt.LockedAmount = amt.Amount + (order.AmountSell + order.Fee)
	// 	err = s.balanceDao.LockFunds(order.UserAddress, order.TokenSell, &amt)

	// 	if err != nil {
	// 		return err
	// 	}

	// } else if order.Type == types.SELL {
	// 	amt := bal.Tokens[order.TokenBuy]
	// 	if amt.Amount < order.AmountBuy+order.Fee {
	// 		return errors.New("Insufficient Balance")
	// 	}
	// 	amt.Amount = amt.Amount - (order.AmountBuy + order.Fee)
	// 	amt.LockedAmount = amt.Amount + (order.AmountBuy + order.Fee)
	// 	err = s.balanceDao.LockFunds(order.UserAddress, order.TokenBuy, &amt)

	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
