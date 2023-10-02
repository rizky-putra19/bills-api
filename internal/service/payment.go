package service

import (
	"errors"
	"time"

	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type Payment struct {
	paymentRepo internal.PaymentRepositoryItf
}

func NewPayment(
	paymentRepo internal.PaymentRepositoryItf,
) *Payment {
	return &Payment{
		paymentRepo: paymentRepo,
	}
}

func (p *Payment) AcceptPayment(referenceNumber string, amount float64) (int, error) {
	payment, err := p.paymentRepo.GetPaymentByReferenceNumber(referenceNumber)
	if err != nil {
		slog.Errorw("failed to accept payment", "err", err.Error())
		return 0, errors.New("payment not found")
	}

	// validate expiry time
	now := time.Now()
	if now.After(*payment.PaymentExpiryTime) {
		slog.Errorw(
			"failed to accept payment",
			"expiry_time", payment.PaymentExpiryTime.String(),
			"now", now.String(),
			"reference_number", referenceNumber,
		)

		return 0, errors.New("payment has been expired")
	}

	// validate amount
	if amount != payment.PaymentAmount {
		slog.Errorw(
			"invalid payment amount",
			"reference_number", referenceNumber,
			"request_amount", payment.PaymentAmount,
			"response_amount", amount,
		)

		return 0, errors.New("invalid payment amount")
	}

	payment.PaymentDate = &now
	payment.Status = constant.PaymentStatusSuccess
	err = p.paymentRepo.UpdatePayment(payment)
	if err != nil {
		slog.Errorw(
			"failed to set payment status",
			"reference_number", referenceNumber,
		)
		return 0, errors.New("invalid payment amount")
	}

	return payment.OrderID, nil
}
