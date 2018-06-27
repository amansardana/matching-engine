package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amansardana/matching-engine/ws"

	"github.com/gorilla/websocket"

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

func (s *PairService) RegisterForOrderBook(conn *websocket.Conn, pairName string) {
	res, err := s.pairDao.GetByName(pairName)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid_Pair_Name",
			"Message": "Invalid Pair Name passed in query Params: " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		conn.Close()
	}
	sKey, bKey := res.GetOrderBookKeys()
	fmt.Printf("\n Sell Key: %s \n Buy Key: %s \n", sKey, bKey)
	// TODO: Get OrderBook from engine

	if err := ws.PairSocketRegister(res.Name, conn); err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER: " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		conn.Close()
	}
	conn.SetCloseHandler(ws.PairSocketCloseHandler(res.Name, conn))

	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			break
		}
	}
}
