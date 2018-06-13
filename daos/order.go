package daos

import (
	"time"

	"github.com/amansardana/matching-engine/app"
	"github.com/amansardana/matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type OrderDao struct {
	collectionName string
	dbName         string
}

func NewOrderDao() *OrderDao {
	return &OrderDao{"orders", app.Config.DBName}
}

func (dao *OrderDao) Create(order *types.Order) (err error) {

	order.ID = bson.NewObjectId()
	order.Status = types.NEW
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, order)
	return
}
