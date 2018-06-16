package services

import (
	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/types"
)

type AddressService struct {
	AddressDao *daos.AddressDao
	balanceDao *daos.BalanceDao
	tokenDao   *daos.TokenDao
}

func NewAddressService(AddressDao *daos.AddressDao, balanceDao *daos.BalanceDao, tokenDao *daos.TokenDao) *AddressService {
	return &AddressService{AddressDao, balanceDao, tokenDao}
}

func (s *AddressService) Create(Address *types.UserAddress) error {
	ua, err := s.GetByAddress(Address.Address)
	if err == nil && ua != nil {
		Address = ua
		return nil
	}
	err = s.AddressDao.Create(Address)
	if err != nil {
		return err
	}
	balService := NewBalanceService(s.balanceDao, s.tokenDao)
	bal := &types.Balance{Address: Address.Address}
	err = balService.Create(bal)
	if err != nil {
		return err
	}
	return err

}

func (s *AddressService) GetByID(id bson.ObjectId) (*types.UserAddress, error) {
	return s.AddressDao.GetByID(id)
}

func (s *AddressService) GetAll() ([]types.UserAddress, error) {
	return s.AddressDao.GetAll()
}
func (s *AddressService) GetByAddress(addr string) (*types.UserAddress, error) {
	return s.AddressDao.GetByAddress(addr)
}
