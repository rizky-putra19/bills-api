package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/converter"
)

func (ctrl *Controller) Authentication(c echo.Context) error {
	var payload dto.AuthPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "invalid payload",
			"stack_trace": err.Error(),
		})
	}

	if payload.Email == "" || payload.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "empty request",
		})
	}

	authResponse, err := ctrl.customerService.Authentication(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "authentication failed",
			"stack_trace": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    authResponse,
		"message": "authentication success",
	})
}

func (ctrl *Controller) CustomerRegistration(c echo.Context) error {
	var payload dto.AuthPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "invalid payload",
			"stack_trace": err.Error(),
		})
	}

	if payload.Email == "" || payload.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "empty request",
		})
	}

	authResponse, err := ctrl.customerService.Registration(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "registration failed",
			"stack_trace": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    authResponse,
		"message": "registration complete",
	})
}

func (ctrl *Controller) FetchProfile(c echo.Context) error {
	custID := c.Get("customer-id")
	cst, err := ctrl.customerService.GetProfile(converter.ToInt(custID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":     "customer data not found",
			"stack_trace": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": cst,
	})
}

func (ctrl *Controller) FetchProduct(c echo.Context) error {
	product, err := ctrl.productService.GetProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":     "product not found",
			"stack_trace": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": product,
	})
}

func (ctrl *Controller) SubmitOrder(c echo.Context) error {
	cfg := ctrl.cfg.Provider
	var payload dto.OrderPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "invalid payload",
			"stack_trace": err.Error(),
		})
	}

	custID := c.Get("customer-id")
	payload.CustomerID = converter.ToInt(custID)

	paymentMetadata, err := ctrl.orderService.SaveOrder(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "order failed",
			"stack_trace": err.Error(),
		})
	}

	// set for redirect payment query
	query := url.Values{}
	query.Add("idPayin", paymentMetadata.ReferenceNumber)
	query.Add("token", paymentMetadata.LFIToken)
	query.Add("amount", converter.ToString(paymentMetadata.Amount))
	query.Add("merchantCode", cfg.MerchantCodeLfi)
	query.Add("currencyCode", constant.IDRCurrencyCode)
	query.Add("callbackUrl", cfg.CallbackUrl)
	query.Add("returnUrl", cfg.ReturnUrl)

	// set url
	u, _ := url.Parse(cfg.PaymentRedirectLfiUrl)
	u.RawQuery = query.Encode()

	fmt.Println(u.String())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"payment_url": u.String(),
			"metadata":    paymentMetadata,
		},
		"message": "waiting payment",
	})
}

func (ctrl *Controller) FetchOrder(c echo.Context) error {
	custID := c.Get("customer-id")
	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "invalid payload",
			"stack_trace": "order id not found",
		})
	}

	orderIDInt := converter.ToInt(orderID)

	orderData, err := ctrl.orderService.FindOrderDetails(orderIDInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "order failed",
			"stack_trace": err.Error(),
		})
	}

	if converter.ToInt(custID) != orderData.CustomerID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":       "forbidden access",
			"stack_trace": "cannot fetch orders belong to another user",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": orderData,
	})
}

func (ctrl *Controller) FetchOrderByCustomer(c echo.Context) error {
	custID := c.Get("customer-id")
	custIDInt := converter.ToInt(custID)

	orderData, err := ctrl.orderService.FindOrderByCustomer(custIDInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "order failed",
			"stack_trace": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": orderData,
	})
}
