package dto

type RequestBillerPayload struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	BillNo   string `json:"billNo"`
	InstCd   string `json:"instCd"`
	Amount   string `json:"amount"`
	RefNo    string `json:"refNo"`
}

type RequestBillerResponsePayload struct {
	Id     string      `json:"id"`
	BillNo string      `json:"billNo"`
	InstCd string      `json:"instCd"`
	Amount string      `json:"amount"`
	RefNo  string      `json:"refNo"`
	Rc     string      `json:"rc"`
	Data   interface{} `json:"data"` // only filled if inquiry is applicable
}
