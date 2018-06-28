package services

import (
	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/types"
)

type TradeService struct {
	tradeDao *daos.TradeDao
}

func NewTradeService(TradeDao *daos.TradeDao) *TradeService {
	return &TradeService{TradeDao}
}

func (t *TradeService) GetByPairName(pairName string) ([]*types.Trade, error) {
	return t.tradeDao.GetByPairName(pairName)
}
