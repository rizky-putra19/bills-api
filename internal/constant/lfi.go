package constant

// Available Payment Processor
const (
	ProcessorPayIn          = "PAYIN"
	ProcessorPayInQuery     = "PAYINQUERY"
	ProcessorPayInCallback  = "PAYINCALLBACK"
	ProcessorPayOut         = "PAYOUT"
	ProcessorPayOutQuery    = "PAYOUTQUERY"
	ProcessorPayOutCallback = "PAYOUTCALLBACK"

	HttpRequestFailed = "failed to send request"
)

const BankAccountNameSimilarityMatchInPercent = 0.4

// Available Payment Method
const (
	PaymentMethodOffline = "OFFLINE"
	PaymentMethodOnline  = "ONLINE"
	PaymentMethodPayout  = "PAYOUT_METHOD"
	PaymentMethodVA      = "VAPAY"

	PayInType  = "PAYIN"
	PayOutType = "PAYOUT"
)

const TenantCode = "LFI"
const OrderIDSeparator = ":::"
const AccountNumberSeparator = ":::"

const IDRCurrencyCode = "IDR"

const (
	PayStatusPending = "PENDING"
	PayStatusSuccess = "SUCCESS"
	PayStatusFailed  = "FAILED"
)

const (
	OrderFailureSourceInternal        = "INTERNAL"
	OrderFailureSourceGatewayRequest  = "GATEWAY_REQUEST"
	OrderFailureSourceGatewayCallback = "GATEWAY_CALLBACK"
)

const (
	OrderFailureCodeSystem                = "SYSTEM"
	OrderFailureCodeBalanceInsufficient   = "BALANCE_INSUFFICIENT"
	OrderFailureCodeAttributeValueInvalid = "ATTRIBUTE_VALUE_INVALID"
	OrderFailureCodeConnectionError       = "CONNECTION_ERROR"
	OrderFailureCodeUnavailablePayout     = "UNAVAILABLE_PAYOUT"
)

const (
	ProviderDataSourceQuery    = "DATA_SOURCE_QUERY"
	ProviderDataSourceCallback = "DATA_SOURCE_CALLBACK"
)

const (
	TransactionTypeInquiry      = "inquiry"
	TransactionTypeVA           = "va"
	TransactionTypeDisbursement = "disbursement"
	ActionTypeRequest           = "Request"
	ActionTypeResponse          = "Response"
)

const (
	LfiMethodPost = "POST"
	LfiMethodGet  = "GET"
)

const (
	LfiUrlPayins           = "payins"
	LfiUrlPayouts          = "payouts"
	LfiUrlAvailablePayment = "!availablePaymentOptions"
)

// ORDER STATUS
// failureSource - one of INTERNAL, GATEWAY_REQUEST, GATEWAY_CALLBACK
// failureCode - one of SYSTEM, ATTRIBUTE_REQUIRED, ATTRIBUTE_INVALID, ATTRIBUTE_VALUE_INVALID,
// ENTITY_NOT_FOUND, ENTITY_DUPLICATED, ENTITY_EXPIRED, CHANNEL_UNAVAILABLE, CHANNEL_LIMIT, MERCHANT_UNAVAILABLE, IP_DENIED,
// BALANCE_INSUFFICIENT, CONNECTION_ERROR, CREDENTIALS_INVALID, SIGNATURE_INVALID, DIFFERENT_AMOUNT_CONFIRMED,
// UNAVAILABLE_PAYOUT, INVALID_REF_NUM, DIFFERENT_AMOUNT
