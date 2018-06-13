package types

import (
	"math"

	"labix.org/v2/mgo/bson"

	"github.com/amansardana/matching-engine/app"
	validation "github.com/go-ozzo/ozzo-validation"
)

type OrderRequest struct {
	TokenBuy         string  `json:"tokenBuy"`
	TokenSell        string  `json:"tokenSell"`
	BuyTokenAddress  string  `json:"buyTokenAddress"`
	SellTokenAddress string  `json:"sellTokenAddress"`
	Type             int     `json:"type" bson:"type"`
	Amount           float64 `json:"amount"`
	Price            float64 `json:"price"`
	Fee              float64 `json:"fee"`
	Signature        string  `json:"signature"`
	PairID           string  `json:"pairID"`
	UserAddress      string  `json:"userAddress"`
}

// Validate validates the OrderRequest fields.
func (m OrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TokenBuy, validation.Required),
		validation.Field(&m.TokenSell, validation.Required),
		validation.Field(&m.BuyTokenAddress, validation.Required),
		validation.Field(&m.SellTokenAddress, validation.Required),
		validation.Field(&m.Type, validation.Required, validation.In(BUY, SELL)),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.UserAddress, validation.Required),
		validation.Field(&m.Signature, validation.Required),
		validation.Field(&m.PairID, validation.Required, validation.NewStringRule(bson.IsObjectIdHex, "Invalid pair id")),
	)
}

// ToOrder converts the OrderRequest to Order
func (m *OrderRequest) ToOrder() (order *Order, err error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	order = &Order{
		TokenBuy:         m.TokenBuy,
		TokenSell:        m.TokenSell,
		BuyTokenAddress:  m.BuyTokenAddress,
		SellTokenAddress: m.SellTokenAddress,
		Type:             OrderType(m.Type),
		Amount:           uint64(m.Amount * math.Pow10(8)),
		Price:            uint64(m.Price * math.Pow10(8)),
		Fee:              uint64(m.Amount * m.Price * (1 + (app.Config.TakeFee / 100))), // amt*price + amt*price*takeFee/100
		PairID:           bson.ObjectIdHex(m.PairID),
		UserAddress:      m.UserAddress,
	}
	return
}
