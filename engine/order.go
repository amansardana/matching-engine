package engine

import "github.com/amansardana/matching-engine/types"

func (e *EngineResource) NewOrder(order *types.Order) (err error) {
	// TODO: Validate if order is valid

	if err = e.orderDao.Create(order); err != nil {
		return
	}

	// TODO: Push order to queue
	
	return
}
