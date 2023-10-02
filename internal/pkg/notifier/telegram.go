package notifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
)

type InquiryMessage struct {
	PayOrderCode  string
	AccountHolder string
	AccountNumber string
	BankCode      string
}

type DisbursementMessage struct {
	PayOrderCode          string
	ProviderTransactionID string
	AccountHolder         string
	AccountNumber         string
	BankCode              string
	Amount                string
}

type VAMessage struct {
	PayOrderCode          string
	ProviderTransactionID string
	AccountHolder         string
	AccountNumber         string
	BankCode              string
	Amount                string
}

type RequestData struct {
	ChatID    int    `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
	Text      string `json:"text"`
}

const InquiryMessageTemplate = `
<b><i>[%v] [GV] Inquiry Account</i></b> %v,

<pre>
Account Number  : %v
Account Holder  : %v
Bank Code       : %v
Create Time	    : %v
</pre>
`

const DisbursementMessageTemplate = `
<b><i>[%v] [LQ] Disbursement</i></b> %v,

<pre>
Provider Tx ID  : %v
Account Number  : %v
Bank Code       : %v
Amount          : %v
Create Time     : %v
</pre>
`

const VAMessageTemplate = `
<b><i>[%v] [LQ] Virtual Account</i></b> %v,

<pre>
Provider Tx ID  : %v
Account Number  : %v
Bank Code       : %v
Amount          : %v
Create Time     : %v
</pre>
`

type Notifier struct {
	client *notify.Notify
}

func buildMessage(msg interface{}, transactionType, action string) string {
	if transactionType == constant.TransactionTypeDisbursement {
		msg := msg.(DisbursementMessage)

		return fmt.Sprintf(
			DisbursementMessageTemplate,
			action,
			msg.PayOrderCode,
			msg.ProviderTransactionID,
			msg.AccountNumber,
			msg.BankCode,
			msg.Amount,
			time.Now().Format(time.RFC3339),
		)
	}

	if transactionType == constant.TransactionTypeVA {
		msg := msg.(VAMessage)

		return fmt.Sprintf(VAMessageTemplate,
			action,
			msg.PayOrderCode,
			msg.ProviderTransactionID,
			msg.AccountNumber,
			msg.BankCode,
			msg.Amount,
			time.Now().Format(time.RFC3339))
	}

	return ""

}

var notifier *Notifier
var once sync.Once

/*
	token	: bot2130842651:AAEzwqQJtIZRmeuLeScHIoeGsZVgov26nqw
	group_id: -668726118 // dev
	group_id: -1001517167280 // prod
*/
func New() {
	once.Do(func() { // <-- atomic, does not allow repeating
		var token = "2130842651:AAEzwqQJtIZRmeuLeScHIoeGsZVgov26nqw"
		var groupID int64 = -1001517167280
		telegramService, _ := telegram.New(token)
		// Passing a telegram chat id as receiver for our messages.
		// Basically where should our message be sent?
		telegramService.AddReceivers(groupID)
		// Create our notifications distributor.
		n := notify.New()
		// Tell our notifier to use the telegram service. You can repeat the above process
		// for as many services as you like and just tell the notifier to use them.
		// Inspired by http middleware used in higher level libraries.
		n.UseServices(telegramService)
		notifier = &Notifier{client: n}
	})
}

func SendMessage(message interface{}, transactionType, action string) {
	msg := buildMessage(message, transactionType, action)
	// Send a test message.
	_ = notifier.client.Send(
		context.Background(),
		"",
		msg,
	)
}
