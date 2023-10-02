package dto

type TokenRequest struct {
	MerchantCode string   `json:"merchantCode"`
	Secret       string   `json:"secret"`
	ValiditySecs int      `json:"validitySecs"`
	Operations   []string `json:"operations"`
}

type TokenRequestResponse struct {
	Token string `json:"token"`
}

type VaRequest struct {
	PaymentRequested PaymentRequestedWrapper `json:"paymentRequested"`
	PaymentMethod    PaymentMethodWrapper    `json:"paymentMethod"`
	CallbackUrl      string                  `json:"callbackUrl"`
	ReturnUrl        string                  `json:"returnUrl"`
}

type VaRequestResponse struct {
	PaymentRequested      PaymentRequestedWrapper      `json:"paymentRequested"`
	PaymentMethodResponse PaymentMethodResponseWrapper `json:"paymentMethodResponse"`
}

type PaymentMethodResponseWrapper struct {
	IdPayin           string         `json:"idPayin"`
	IdPayment         string         `json:"idPayment"`
	Account           AccountWrapper `json:"account"`
	Money             MoneyWrapper   `json:"money"`
	Reference         string         `json:"reference"`
	ReturnUrl         string         `json:"returnUrl"`
	AcceptedAt        string         `json:"acceptedAt"`
	ExpireAt          string         `json:"expireAt"`
	PaymentMethodCode string         `json:"paymentMethodCode"`
}

type PaymentMethodWrapper struct {
	PaymentMethodCode string         `json:"paymentMethodCode"`
	Account           AccountWrapper `json:"account"`
}

type AccountWrapper struct {
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber,omitempty"`
	BankCode      string `json:"bankCode,omitempty"`
	BankName      string `json:"bankName,omitempty"`
	BankBranch    string `json:"bankBranch,omitempty"`
	BankCity      string `json:"bankCity,omitempty"`
	BankProvince  string `json:"bankProvince,omitempty"`
}

type PaymentRequestedWrapper struct {
	Money       MoneyWrapper `json:"money"`
	SegmentCode string       `json:"segmentCode"`
}

type PaymentInstructionListWrapper struct {
	BankService string `json:"bankingService"`
	CompanyCode string `json:"companyCode"`
}

type MoneyWrapper struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

type SyncStatusResponse struct {
	PaymentRequest        PaymentRequestedWrapper      `json:"paymentRequested"`
	PaymentMethodResponse PaymentMethodResponseWrapper `json:"paymentMethodResponse"`
	Fee                   FeeWrapper                   `json:"fee"`
	Process               ProcessWrapper               `json:"process"`
}

type FeeWrapper struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type ProcessWrapper struct {
	Status          string               `json:"status"`
	FailureReasons  FailureReasonWrapper `json:"failureReasons"`
	CreatedAt       string               `json:"createdAt"`
	ProcessedAt     string               `json:"processedAt"`
	IsTest          bool                 `json:"isTest"`
	ProcessorStatus string               `json:"processorStatus"`
}

type FailureReasonWrapper struct {
	FailedAt          string `json:"failedAt"`
	FailureReasonCode string `json:"failureReasonCode"`
	Message           string `json:"message"`
}

type LFIPayinCallbackPayload struct {
	IdEntity string `json:"idEntity"`
}
