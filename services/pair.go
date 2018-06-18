package services

import (
	"errors"

	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/daos"
	aerrors "github.com/amansardana/matching-engine/errors"
	"github.com/amansardana/matching-engine/types"
)

type PairService struct {
	pairDao  *daos.PairDao
	tokenDao *daos.TokenDao
}

func NewPairService(pairDao *daos.PairDao, tokenDao *daos.TokenDao) *PairService {
	return &PairService{pairDao, tokenDao}
}

func (s *PairService) Create(pair *types.Pair) error {
	bt, err := s.tokenDao.GetByID(pair.BuyToken)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.BuyToken.Hex() + " doesn't exists")})
	}
	st, err := s.tokenDao.GetByID(pair.SellToken)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.SellToken.Hex() + " doesn't exists")})
	}
	pair.SellTokenSymbol = st.Symbol
	pair.SellTokenAddress = st.ContractAddress
	pair.BuyTokenSymbol = bt.Symbol
	pair.BuyTokenAddress = bt.ContractAddress

	err = s.pairDao.Create(pair)
	return err

}

func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}
