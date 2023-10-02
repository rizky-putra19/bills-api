package dto

type JSONResponse struct {
	NextPayInUrl    string                 `json:"nextPayinUrl"`
	NextParams      map[string]interface{} `json:"nextParams"`
	RawResponseStr  string                 `json:"rawResponseStr"`
	ResponseMessage string                 `json:"responseMessage"`
	ResponseCode    string                 `json:"responseCode"`
}
