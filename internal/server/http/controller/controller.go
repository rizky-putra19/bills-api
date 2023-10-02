package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/converter"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
)

type Controller struct {
	cfg                config.Schema
	productService     internal.ProductServiceItf
	paymentService     internal.PaymentServiceItf
	orderService       internal.OrderServiceItf
	fulfillmentService internal.FulfillmentServiceItf
	customerService    internal.CustomerServiceItf
}

func NewController(
	cfg config.Schema,
	product internal.ProductServiceItf,
	payment internal.PaymentServiceItf,
	order internal.OrderServiceItf,
	fulfillment internal.FulfillmentServiceItf,
	customer internal.CustomerServiceItf,
) *Controller {
	return &Controller{
		cfg:                cfg,
		productService:     product,
		paymentService:     payment,
		orderService:       order,
		fulfillmentService: fulfillment,
		customerService:    customer,
	}
}

/** will be use for kubernetes deployment

func (ctrl *Controller) ReturnAlive(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"status":    200,
		"app_name":  "PROVIDER_NAME_provider",
		"message":   "service live and well",
	})
}

func (ctrl *Controller) ReturnReady(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"status":    200,
		"app_name":  "PROVIDER_NAME_provider",
		"message":   "app ready to served",
	})
}
*/

func (ctrl *Controller) ReturnOK(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (ctrl *Controller) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing Authorization header"})
		}

		token, err := request.ParseFromRequest(c.Request(), request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(ctrl.cfg.App.JWTSecret), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		// Verify the token and check the expiration time
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		currentTime := time.Now().Unix()
		if currentTime > converter.ToInt64(claims["exp"]) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token has expired"})
		}

		c.Set("customer-id", claims["id"])
		// slog.Infow("claimed token", "data", claims)

		return next(c)
	}
}
