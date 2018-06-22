package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/amansardana/matching-engine/engine"
	"github.com/amansardana/matching-engine/services"
	"github.com/amansardana/matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
)

type orderEndpoint struct {
	orderService *services.OrderService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeOrder sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(rg *routing.RouteGroup, orderService *services.OrderService, e *engine.EngineResource) {
	r := &orderEndpoint{orderService}
	// rg.Get("/orders/<id>", r.get)
	// rg.Get("/orders", r.query)
	rg.Post("/orders", r.create)
	rg.Any("/orders/ws", r.ws)
	e.SubscribeEngineResponse(r.engineResponse)
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
	err = r.orderService.Create(order)
	if err != nil {
		return err
	}

	return c.Write(order)
}
func (r *orderEndpoint) ws(c *routing.Context) error {
	conn, err := upgrader.Upgrade(c.Response, c.Request, nil)
	if err != nil {
		log.Println("==>" + err.Error())
		return nil
	}
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println("<==>" + err.Error())
		return nil
	}
	if err := conn.WriteMessage(messageType, p); err != nil {
		log.Println("<<==>>" + err.Error())
		return nil
	}
	return nil
}

func (r *orderEndpoint) engineResponse(er *engine.EngineResponse) error {
	b, _ := json.Marshal(er)
	fmt.Printf("\n======> \n%s\n <======\n", b)
	return nil
}
