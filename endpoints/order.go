package endpoints

import (
	"fmt"

	"github.com/amansardana/matching-engine/services"
	"github.com/amansardana/matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
)

type orderEndpoint struct {
	orderService *services.OrderService
}

// ServeOrder sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(rg *routing.RouteGroup, orderService *services.OrderService) {
	r := &orderEndpoint{orderService}
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
	order, err := model.ToOrder()
	if err != nil {
		return err
	}
	fmt.Println(order)
	err = r.orderService.Create(order)
	if err != nil {
		return err
	}

	return c.Write(order)
}
