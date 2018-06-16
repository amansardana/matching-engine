package services

import (
	"math"

	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/types"
)

type BalanceService struct {
	balanceDao *daos.BalanceDao
	tokenDao   *daos.TokenDao
}

func NewBalanceService(balanceDao *daos.BalanceDao, tokenDao *daos.TokenDao) *BalanceService {
	return &BalanceService{balanceDao, tokenDao}
}

func (s *BalanceService) Create(balance *types.Balance) error {
	tb := make(map[string]types.TokenBalance)
	tokens, err := s.tokenDao.GetAll()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		tb[token.Symbol] = types.TokenBalance{Amount: uint64(1000 * math.Pow10(8)), LockedAmount: 0, TokenID: token.ID}
	}
	balance.Tokens = tb
	err = s.balanceDao.Create(balance)
	return err

}

func (s *BalanceService) GetByID(id bson.ObjectId) (*types.Balance, error) {
	return s.balanceDao.GetByID(id)
}
func (s *BalanceService) GetByAddress(addr string) (*types.Balance, error) {
	return s.balanceDao.GetByAddress(addr)
}
func (s *BalanceService) GetAll() ([]types.Balance, error) {
	return s.balanceDao.GetAll()
}
