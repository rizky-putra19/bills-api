package prismalink

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/converter"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/helper"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type prismalink struct {
	config     config.Provider
	httpClient *http.Client
}

func New(cfg config.Provider) *prismalink {
	return &prismalink{
		config: cfg,
		httpClient: &http.Client{
			Timeout: constant.NinetySecond,
		},
	}
}

func (prismalink *prismalink) RequestBiller(payload dto.RequestPayloadPrismalink, reqType string) (dto.RequestBillerResponsePayload, error) {
	cfg := prismalink.config
	var schema dto.RequestBillerResponsePayload

	// password : MD5(merchantId+IdTransaction+clientKey+amount+billNo)
	stringToEncode := cfg.MerchantIdPrismalink +
		payload.IdTransaction +
		cfg.ClientKeyPrismalink +
		payload.Amount +
		payload.BillNo

	password := helper.EncodeToMd5(stringToEncode)

	// request payload
	requestPayload := dto.RequestBillerPayload{
		Id:       cfg.MerchantIdPrismalink,
		Password: password,
		BillNo:   payload.BillNo,
		InstCd:   payload.InstCd,
		Amount:   payload.Amount,
		RefNo:    payload.IdTransaction,
	}

	// request payload for Log
	logPayload := requestPayload
	logPayload.Password = constant.PrismalinkPassCensored

	// transform to payload json
	payloadJson, _ := json.Marshal(requestPayload)

	// http request
	r, err := http.NewRequest(http.MethodPost, cfg.RequestPrismalinkURL, bytes.NewBuffer(payloadJson)) // json payload
	if err != nil {
		return schema, err
	}

	// start block code for inquiry request
	if reqType == constant.PrismalinkInquiryRequest {
		// password : MD5(merchantId+IdTransaction+clientKey+amount+billNo)
		stringToEncode = cfg.MerchantIdPrismalink +
			payload.IdTransaction +
			cfg.ClientKeyPrismalink +
			constant.PrismalinkInquiryValue +
			payload.BillNo

		password = helper.EncodeToMd5(stringToEncode)

		// request payload
		requestPayload = dto.RequestBillerPayload{
			Id:       cfg.MerchantIdPrismalink,
			Password: password,
			BillNo:   payload.BillNo,
			InstCd:   payload.InstCd,
			Amount:   constant.PrismalinkInquiryValue,
			RefNo:    payload.IdTransaction,
		}

		// transform payload to json
		payloadJson, _ = json.Marshal(requestPayload)

		// set new body payload to http.Request for inquiry
		r.Body = ioutil.NopCloser(bytes.NewBuffer(payloadJson))

		// request payload for log
		logPayload = requestPayload
		logPayload.Password = constant.PrismalinkPassCensored
	}
	// end of block code inquiry request

	// logging payload request
	slog.Infof("PRISMALINK %v [%v-request] payload: %v", payload.IdTransaction, strings.ToLower(reqType), converter.ToString(logPayload))

	// set header
	r.Header.Add("Content-Type", "application/json")
	r.Close = true

	response, err := prismalink.httpClient.Do(r)
	if err != nil {
		return schema, err
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println("failed to close response body, could lead to memory leak")
		}
	}()

	// read response body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		slog.Infof("PRISMALINK %v [%v-request] got error failed to read response %v", payload.IdTransaction, strings.ToLower(reqType), string(contents))
		return schema, err
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("PRISMALINK %v [%v-request] got error status not ok %v", payload.IdTransaction, strings.ToLower(reqType), string(contents))
		return schema, errors.New("status not ok")
	}

	// convert response to struct
	err = json.Unmarshal(contents, &schema)
	if err != nil {
		slog.Infof("PRISMALINK %v [%v-request] got error failed unmarshal response %v", payload.IdTransaction, strings.ToLower(reqType), string(contents))
		return schema, err
	}

	slog.Infof("PRISMALINK %v [%v-request] response data: %v", payload.IdTransaction, strings.ToLower(reqType), converter.ToString(schema))

	return schema, nil
}
