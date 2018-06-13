package engine

import (
	"errors"

	"github.com/amansardana/matching-engine/daos"
)

type EngineResource struct {
	orderDao *daos.OrderDao
}

var Engine *EngineResource

func InitEngine(orderDao *daos.OrderDao) (err error) {
	if Engine == nil {
		if orderDao == nil {
			return errors.New("Need pointer to struct of type daos.OrderDao")
		}
		Engine = &EngineResource{orderDao}
	}
	return
}
