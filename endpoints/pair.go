package endpoints

import (
	"github.com/amansardana/matching-engine/errors"
	"github.com/amansardana/matching-engine/pairs"
	"github.com/amansardana/matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
	"labix.org/v2/mgo/bson"
)

type pairEndpoint struct {
	pairService *pairs.PairService
}

// ServePair sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(rg *routing.RouteGroup, pairService *pairs.PairService) {
	r := &pairEndpoint{pairService}
	rg.Get("/pairs/<id>", r.get)
	rg.Get("/pairs", r.query)
	rg.Post("/pairs", r.create)
}

func (r *pairEndpoint) create(c *routing.Context) error {
	var model types.Pair
	if err := c.Read(&model); err != nil {
		return err
	}
	if err := model.Validate(); err != nil {
		return err
	}
	err := r.pairService.Create(&model)
	if err != nil {
		return err
	}

	return c.Write(model)
}

func (r *pairEndpoint) query(c *routing.Context) error {

	response, err := r.pairService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *pairEndpoint) get(c *routing.Context) error {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return errors.NewAPIError(400, "INVALID_ID", nil)
	}
	response, err := r.pairService.GetByID(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}
