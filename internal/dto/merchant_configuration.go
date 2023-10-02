package dto

type MerchantConfig struct {
	TenantCode        string `json:"tenantCode"`
	MerchantCode      string `json:"merchantCode"`
	MerchantName      string `json:"merchantName"`
	PayType           string `json:"payType"`
	PayOrderId        string `json:"payOrderId"`
	MchOrderId        string `json:"mchOrderId"`
	BankCode          string `json:"bankCode"`
	Amount            string `json:"amount"`
	ProviderAccountId string `json:"providerAccountId"`
}
