package main

import (
	"fmt"
	"net/http"

	"github.com/amansardana/matching-engine/pairs"
	"github.com/amansardana/matching-engine/tokens"

	"github.com/amansardana/matching-engine/endpoints"

	"github.com/amansardana/matching-engine/engine"

	"github.com/amansardana/matching-engine/app"
	"github.com/amansardana/matching-engine/daos"
	"github.com/amansardana/matching-engine/errors"
	"github.com/Sirupsen/logrus"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
)

func main() {
	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	// load error messages
	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// create the logger
	logger := logrus.New()

	// connect to the database
	if err := daos.InitSession(); err != nil {
		panic(err)
	}
	// db.LogFunc = logger.Infof

	// wire up API routing
	http.Handle("/", buildRouter(logger))

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func buildRouter(logger *logrus.Logger) *routing.Router {
	router := routing.New()

	router.To("GET,HEAD", "/ping", func(c *routing.Context) error {
		c.Abort() // skip all other middlewares/handlers
		return c.Write("OK " + app.Version)
	})

	router.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	rg := router.Group("")

	// rg.Post("/auth", apis.Auth(app.Config.JWTSigningKey))
	// rg.Use(auth.JWT(app.Config.JWTVerificationKey, auth.JWTOptions{
	// 	SigningMethod: app.Config.JWTSigningMethod,
	// 	TokenHandler:  apis.JWTHandler,
	// }))

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()

	// get services for injection
	tokenService := tokens.NewTokenService(tokenDao)
	pairService := pairs.NewPairService(pairDao, tokenDao)

	// instantiate engine
	if err := engine.InitEngine(orderDao); err != nil {
		panic(err)
	}

	endpoints.ServeTokenResource(rg, tokenService)
	endpoints.ServePairResource(rg, pairService)
	endpoints.ServeOrderResource(rg)
	// artistDAO := daos.NewArtistDAO()
	// apis.ServeArtistResource(rg, services.NewArtistService(artistDAO))

	// wire up more resource APIs here

	return router
}
