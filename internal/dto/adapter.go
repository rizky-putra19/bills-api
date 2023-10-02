package dto

type VAPayload struct {
	BankCode     string
	ExpiryDate   string
	PayOrderCode string
	Amount       int
}

type InquiryPayload struct {
	BankCode      string
	AccountNumber string
	Amount        int
	PayOrderCode  string // will translate to partner_reff
}

type TransferPayload struct {
	BankCode      string
	AccountNumber string
	Amount        int
	PayOrderCode  string // will translate to partner_reff
	InquiryReff   string // will get this value from inquiry response
}

type PayRequest struct {
	BankCode      string
	Id            string
	AccountNumber string
	AccountName   string
	Amount        int
	Email         string
}

type RequestPayloadPrismalink struct {
	IdTransaction string // Transaction ID from database
	BillNo        string // filled by bill number (ex: prepaid phone filled by phone number)
	InstCd        string // institution code of biller
	Amount        string // amount to pay
}
