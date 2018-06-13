package daos

import (
	"time"

	"github.com/amansardana/matching-engine/app"
	"github.com/amansardana/matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type PairDao struct {
	collectionName string
	dbName         string
}

func NewPairDao() *PairDao {
	return &PairDao{"pairs", app.Config.DBName}
}

func (dao *PairDao) Create(pair *types.Pair) (err error) {

	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, pair)
	return
}

func (dao *PairDao) GetAll() (response []types.Pair, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

func (dao *PairDao) GetByID(id bson.ObjectId) (response *types.Pair, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
