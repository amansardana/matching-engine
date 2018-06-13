package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"labix.org/v2/mgo/bson"
)

type Pair struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	Code          string        `json:"code" bson:"code"`
	Name          string        `json:"name" bson:"name"`
	BuyToken      bson.ObjectId `json:"buyToken" bson:"buyToken"`
	BuyTokenCode  string        `json:"buyTokenCode" bson:"buyTokenCode"`
	SellToken     bson.ObjectId `json:"sellToken" bson:"sellToken"`
	SellTokenCode string        `json:"sellTokenCode" bson:"sellTokenCode"`

	MakerFee float64 `json:"makerFee" bson:"makerFee"`
	TakerFee float64 `json:"takerFee" bson:"takerFee"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (t Pair) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Code, validation.Required),
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.BuyToken, validation.Required),
		validation.Field(&t.BuyTokenCode, validation.Required),
		validation.Field(&t.SellToken, validation.Required),
		validation.Field(&t.SellTokenCode, validation.Required),
	)
}
