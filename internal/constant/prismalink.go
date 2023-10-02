package constant

const PrismalinkInquiryValue = "0"
const PrismalinkPassCensored = "********"

const (
	PrismalinkInquiryRequest = "INQUIRY"
	PrismalinkPaymentRequest = "PAYMENT"
)

const (
	PLIndosatInstitutionCode   = "indosatIm3Prepaid"
	PLTelkomsepInstitutionCode = "telkomselSimpatiPrepaid"
	PLXLInstitutionCode        = "xlBebasPrepaid"
	PLTriInstitutionCode       = "3Prepaid"
)

const (
	PrismalinkBillerResponseSuccess           = "00"
	PrismalinkBillerResponseUnknownBillNo     = "14"
	PrismalinkBillerResponseWrongUserPassword = "27"
	PrismalinkBillerResponseTransactionError  = "20"
)

var PrismalinkPrePaidProduct = map[string]string{
	"25K":  "25000",
	"50K":  "50000",
	"100K": "100000",
	"150K": "150000",
}
