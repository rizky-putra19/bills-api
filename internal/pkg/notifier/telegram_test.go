package notifier

import (
	"testing"

	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
)

func TestSendMessage(t *testing.T) {
	New()
	type args struct {
		message         interface{}
		transactionType string
		action          string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success transaction",
			args: args{
				message: DisbursementMessage{
					PayOrderCode:          "xxx",
					ProviderTransactionID: "00",
					AccountHolder:         "wardana",
					AccountNumber:         "10000",
					BankCode:              "bca",
					Amount:                "10000",
				},
				transactionType: constant.TransactionTypeDisbursement,
				action:          constant.ActionTypeRequest,
			},
		},
		{
			name: "success transaction inquiry",
			args: args{
				message: InquiryMessage{
					PayOrderCode:  "xxx",
					AccountHolder: "wardana",
					AccountNumber: "10000",
					BankCode:      "bca",
				},
				transactionType: constant.TransactionTypeInquiry,
				action:          constant.ActionTypeRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMessage(tt.args.message, tt.args.transactionType, tt.args.action)
		})
	}
}
