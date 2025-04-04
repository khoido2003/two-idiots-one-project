package routers

import (
	"server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type RouterContext struct {
	v1Router       *chi.Mux
	handlerContext handlers.HandlerContext
}

func NewRouter(handlerContext handlers.HandlerContext) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Test route
	router.Get("/ping", handlerContext.Ping)
	router.Get("/error", handlerContext.HandleError)

	// Create routes handler
	routerContext := RouterContext{
		v1Router:       chi.NewRouter(),
		handlerContext: handlerContext,
	}

	// Routes
	routerContext.CheckHealth()
	routerContext.ProductsRoute()
	routerContext.CategoriesRoute()
    routerContext.UsersRoute()
    routerContext.CartsRoute()

	// Connnect v1  routes to the main routes
	router.Mount("/api/v1", routerContext.v1Router)

	return router

}
