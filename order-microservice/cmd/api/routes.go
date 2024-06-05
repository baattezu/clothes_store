package main

import (
	"net/http"
	"order-microservice/internal/auth"
	"order-microservice/internal/cart"
	"order-microservice/internal/order"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	cartRepo := cart.NewRepository(app.db)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	orderRepo := order.NewRepository(app.db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	router.Handler(http.MethodGet, "/cart", wrapHandler(auth.AuthMiddleware(http.HandlerFunc(cartHandler.HandleCart))))
	router.Handler(http.MethodPost, "/cart", wrapHandler(auth.AuthMiddleware(http.HandlerFunc(cartHandler.HandleCart))))
	router.Handler(http.MethodPost, "/order", wrapHandler(auth.AuthMiddleware(http.HandlerFunc(orderHandler.HandleOrder))))
	router.Handler(http.MethodGet, "/v1/healthcheck", http.HandlerFunc(app.healthcheckHandler))

	return app.recoverPanic(app.rateLimit(router))
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
