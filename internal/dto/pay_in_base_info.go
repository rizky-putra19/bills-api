package dto

type PayInBaseInfo struct {
	PayOrderID  string `json:"PayOrderId"`
	Amount      string `json:"amount"`
	BankCode    string `json:"bankCode"`
	DepositName string `json:"depositName"`
	ClientIP    string `json:"clientIp"`
	ReturnUrl   string `json:"returnUrl"`
}
