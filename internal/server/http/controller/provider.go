package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

func (ctrl *Controller) PaymentCallback(c echo.Context) error {
	var payload dto.LFIPayinCallbackPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":       "invalid payload",
			"stack_trace": err.Error(),
		})
	}
	slog.Infof("LFI %v http-request /lfi/payment-callback [start]", payload.IdEntity)

	// proceed request biller
	billerResponse, err := ctrl.fulfillmentService.OrderFulfillmentPrismalinkProvider(payload.IdEntity)
	if err != nil {
		slog.Infof("LFI %v http-request /lfi/payment-callback [end] [failed] request to biller", payload.IdEntity)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":       "failed request to biller",
			"stack_trace": err.Error(),
		})
	}

	slog.Infof("LFI %v http-request /lfi/payment-callback [end] [success] request to biller: %v", payload.IdEntity, billerResponse)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    billerResponse,
		"message": "successfully request to biller",
	})
}
