package dto

type Response struct {
	Status               int    `json:"status"`
	IsContinueToQuery    bool   `json:"isContinueToQuery"`
	ResponseFromProvider string `json:"responseFromProvider"`
}
