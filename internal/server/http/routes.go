package http

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/lokalpay-dev/digital-goods/internal/server/http/controller"
)

func RegisterRoutes(e *echo.Echo, ctrl *controller.Controller) {
	e.GET("/", ctrl.ReturnOK)

	// restricted route
	r := e.Group("/backoffice/v1")
	r.GET("/", ctrl.ReturnOK)
	r.GET("/product", ctrl.ReturnOK)
	r.GET("/order", ctrl.ReturnOK)
	r.GET("/customer", ctrl.ReturnOK)
	r.GET("/customer/:id", ctrl.ReturnOK)

	// public endpoint
	pub := e.Group("/v1")
	pub.POST("/register", ctrl.CustomerRegistration)
	pub.POST("/login", ctrl.Authentication)
	pub.GET("/profile", ctrl.AuthMiddleware(ctrl.FetchProfile))
	pub.GET("/list-product", ctrl.FetchProduct)
	pub.POST("/order", ctrl.AuthMiddleware(ctrl.SubmitOrder))
	pub.GET("/order/:id", ctrl.AuthMiddleware(ctrl.FetchOrder))
	pub.GET("/order-history", ctrl.AuthMiddleware(ctrl.FetchOrderByCustomer))

	// provider callback
	pvdr := e.Group("/provider/v1")
	pvdr.POST("/lfi/payment-callback", ctrl.PaymentCallback)
}
