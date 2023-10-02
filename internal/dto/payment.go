package dto

import "time"

type PaymentPayload struct {
	Timestamp       string                 `json:"timestamp"`
	ProcessPayOrder ProcessPayOrderWrapper `json:"processPayOrder"`
}

type GetPaymentPayload struct {
	Timestamp     string               `json:"timestamp"`
	QueryPayOrder QueryPayOrderWrapper `json:"queryPayOrder"`
}

type QueryPayOrderWrapper struct {
	PayOrder         PayOrderWrapper `json:"payOrder"`
	QueryStartedAt   string          `json:"queryStartedAt"`
	LastAttemptAt    string          `json:"lastAttemptAt"`
	AttemptNumber    int             `json:"attemptNumber"`
	ExpireAt         string          `json:"expireAt"`
	InterfaceSetting string          `json:"interfaceSetting"`
	ReturnURL        string          `json:"returnUrl"`
}

type ProcessPayOrderWrapper struct {
	ApiVersion int             `json:"apiVersion"`
	PayOrder   PayOrderWrapper `json:"payOrder"`
}

type PayOrderWrapper struct {
	ID                  PayOrderID                `json:"idPayOrder" dynamo:"idPayOrder"`
	PayMethod           string                    `json:"payMethod" dynamo:"payMethod"`
	PayType             string                    `json:"payType"`
	Amount              string                    `json:"amount" dynamo:"amount"`
	Currency            string                    `json:"currency" dynamo:"currency"`
	Creditor            Creditor                  `json:"creditor" dynamo:"creditor"`
	Debitor             Debitor                   `json:"debitor" dynamo:"debitor"`
	Provider            Provider                  `json:"provider" dynamo:"provider"`
	ExpirationInSeconds int                       `json:"expirationInSeconds" dynamo:"expirationInSeconds"`
	References          PayOrderReferencesWrapper `json:"references" dynamo:"references"`
	InterfaceSetting    string                    `json:"interfaceSetting" dynamo:"interfaceSetting"`
	ReturnURL           string                    `json:"returnUrl" dynamo:"returnUrl"`
}

type PayOrderID struct {
	TenantCode   string `json:"tenantCode" dynamo:"tenantCode"`
	PayOrderCode string `json:"payOrderCode" dynamo:"payOrderCode"`
}

type Creditor struct {
	AccountName      string `json:"accountName"`
	AccountNumber    string `json:"accountNumber"`
	BankCode         string `json:"bankCode"`
	ReferenceNumber  string `json:"referenceNum"`
	BankName         string `json:"bankName"`
	BankCodeExternal string `json:"bankCodeExternal"`
	Email            string `json:"email"`
	AccountType      string `json:"accountType"`
}

type Debitor struct {
	BankCode        string `json:"bankCode"`
	BankName        string `json:"bankName"`
	ReferenceNumber string `json:"referenceNum"`
	AccountType     string `json:"accountType"`
}

type Provider struct {
	ProviderCode        string `json:"providerCode" dynamo:"providerCode"`
	ProviderName        string `json:"providerName" dynamo:"providerName"`
	ProviderAccountCode string `json:"providerAccountCode" dynamo:"providerAccountCode"`
}

type PayOrderFailedWrapper struct {
	Timestamp string         `json:"timestamp"`
	Order     PayOrderFailed `json:"payOrderFailed"`
}

type PayOrderFailed struct {
	ID                 PayOrderID                `json:"idPayOrder"`
	FailedAt           time.Time                 `json:"failedAt"`
	FailureSource      string                    `json:"failureSource"`
	FailureCode        string                    `json:"failureCode"`
	FailedMessage      string                    `json:"failedMessage"`
	References         PayOrderReferencesWrapper `json:"references"`
	ResponseCode       string                    `json:"responseCode"`
	ProviderDataSource string                    `json:"providerDataSource"`
}

type PayOrderAcceptedByProviderWrapper struct {
	Timestamp                  string                     `json:"timestamp"`
	PayOrderAcceptedByProvider PayOrderAcceptedByProvider `json:"payOrderAcceptedByProvider"`
}

type PayOrderAcceptedByProvider struct {
	ID                  PayOrderID                `json:"idPayOrder"`
	AcceptedAt          time.Time                 `json:"acceptedAt"`
	ProviderResponse    string                    `json:"providerResponse"`
	Creditor            Creditor                  `json:"creditor"`
	Debitor             Debitor                   `json:"debitor"`
	PaymentSpecificData *PaymentSpecificData      `json:"paymentSpecificData"`
	References          PayOrderReferencesWrapper `json:"references"`
	InterfaceSetting    string                    `json:"interfaceSetting" dynamo:"interfaceSetting"`
}
type PayOrderReferencesWrapper struct {
	IDPaymentProcess string `json:"ID_PAYMENT_PROCESS" dynamo:"ID_PAYMENT_PROCESS"`
}

// empty struct
type PaymentSpecificData struct {
	QrCode string `json:"QR_CODE"`
}

type PayOrderSucceededWrapper struct {
	Timestamp string            `json:"timestamp"`
	Order     PayOrderSucceeded `json:"payOrderSucceeded"`
}

type PayOrderStatusQueryAttemptWrapper struct {
	Timestamp string                     `json:"timestamp"`
	Order     PayOrderStatusQueryAttempt `json:"payOrderStatusQueryAttempt"`
}

type PayOrderStatusQueryAttempt struct {
	ID         PayOrderID                `json:"idPayOrder"`
	AttemptAt  time.Time                 `json:"attemptAt"`
	References PayOrderReferencesWrapper `json:"references"`
}

type PayOrderSucceeded struct {
	ID                 PayOrderID `json:"idPayOrder"`
	SucceededAt        time.Time  `json:"succeededAt"`
	ForMoney           ForMoney   `json:"forMoney"`
	ProviderDataSource string     `json:"providerDataSource"`
}

type ForMoney struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type PayOrderFailedData struct {
	PayOrderID    PayOrderID
	References    PayOrderReferencesWrapper
	FailureCode   string
	FailedMessage string
	Status        string
}

/*
{
    "timestamp":"1618494921176",
    "payOrderFailed":{
        "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
        "failedAt":"2021-04-15T13:55:21.176256Z",
        "failureSource":"GATEWAY_CALLBACK",
        "failureCode":"SYSTEM",
        "failedMessage":"Payment failed via callback",
        "responseCode":"EXPIRED"
    }
}
*/

/*
	{
  "timestamp": "1627634612192",
  "processPayOrder": {
    "apiVersion": 2,
    "payOrder": {
      "idPayOrder": {
        "tenantCode": "CAR",
        "payOrderCode": "CAR-IDR-P0-9becc193-f070-42d0-915a-34fc398949dd"
      },
      "payMethod": "PAYOUT_METHOD",
      "payType": "PAYOUT",
      "amount": "10000",
      "currency": "IDR",
      "merchantCode": "CAR21-N2UH58T5LW",
      "merchantName": "LKR Test 1",
      "aggrPayChannelCode": "CAR21-N2UH58T5LW_IDR_PAYOUT_METHOD",
      "payChannelCode": "INTRAJSA-PAYOUT",
      "payChannelVertical": "PAY_VERTICAL_UNKNOWN",
      "expirationInSeconds": 7200,
      "idPayOrderFromMerchant": "TST-246",
      "creditor": {
        "accountName": "Test account",
        "accountNumber": "11223344",
        "bankCode": "IDR_009",
        "bankName": "Bank BNI",
        "email": "customer@test.com",
        "accountType": "UNSET"
      },
      "debitor": {
        "bankCode": "IDR_009",
        "bankName": "Bank BNI",
        "referenceNum": "34fc398949dd",
        "accountType": "UNSET"
      },
      "provider": {
        "providerCode": "INTRAJASA",
        "providerName": "Intrajasa",
        "providerAccountCode": "INTRAJASA-IDR-INTRAJASA-IDR"
      },
      "merchantFee": {
        "calcMethod": "PERCENT",
        "value": "2.0"
      },
      "platformFee": {
        "calcMethod": "AMOUNT",
        "value": "1.0"
      },
      "segment": "SEGMENT_CODE_UNKNOWN",
      "clientIp": "185.63.98.56",
      "isTest": true,
      "returnUrl": "",
      "callbackUrl": "https://enh1vuvwvglylzh.m.pipedream.net",
      "redirectUrl": "",
      "paymentSpecificData": {},
      "references": {
        "ID_PAYMENT_PROCESS": "akka://router@10.0.2.62:2552/user/$b/d3a1c45b-64d2-408b-bfc6-54a7937886b6#1899165744"
      }
    },
    "acceptedAt": "2021-07-30T08:43:32.121172Z",
    "interfaceSetting": "{}",
    "returnUrl": ""
  }
}
*/
