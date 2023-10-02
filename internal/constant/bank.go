package constant

// WhiteListBank is list of enabled bank to create VA
var WhiteListBank = []string{
	"IDR_014", // Bank BCA
	"IDR_002", // Bank BRI
	"IDR_022", // Bank CIMB
	"IDR_009", // Bank BNI
	"IDR_008", // Bank MANDIRI
	"IDR_523", // Bank SAHABAT SAMPOERNA
	"IDR_016", // Bank Maybank
}

var TransformToVAPrefix = map[string]int{
	"IDR_008": 89039099, // 5 digits BIN + 3 digits merchant identifier - mandiri
	"IDR_013": 850899,   // 4 digits BIN + 2 digits merchant identifier - permata
	"IDR_014": 13887,    // 5 digits BIN - bca
	"IDR_009": 8919099,  // 4 digits BIN + 3 digits merchant identifier - BNI
	"IDR_002": 40419099, // 5 digits BIN + 3 digits merchant identifier - BRI
}
