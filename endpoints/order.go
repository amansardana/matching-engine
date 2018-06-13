package endpoints

import (
	"github.com/amansardana/matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
)

type orderEndpoint struct{}

// ServeOrder sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(rg *routing.RouteGroup) {
	r := &orderEndpoint{}
	// rg.Get("/orders/<id>", r.get)
	// rg.Get("/orders", r.query)
	rg.Post("/orders", r.create)
}

func (r *orderEndpoint) create(c *routing.Context) error {
	var model types.OrderRequest
	if err := c.Read(&model); err != nil {
		return err
	}
	if err := model.Validate(); err != nil {
		return err
	}
	// response, err := r.service.Create(app.GetRequestScope(c), &model)
	// if err != nil {
	// 	return err
	// }

	return c.Write("")
}
