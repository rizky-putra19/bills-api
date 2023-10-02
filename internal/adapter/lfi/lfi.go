package lfi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/converter"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type lfi struct {
	config     config.Provider
	httpClient *http.Client
}

func New(cfg config.Provider) *lfi {
	return &lfi{
		config: cfg,
		httpClient: &http.Client{
			Timeout: constant.NinetySecond,
		},
	}
}

func (lfi *lfi) RequestAuthToken(id string) (string, error) {
	cfg := lfi.config
	oneHour := constant.OneHour
	// operationPostAvailablePayment for ability to the choose payment method / operator by user
	operationPostAvailablePayment := fmt.Sprintf("%v /%v/%v", constant.LfiMethodPost, constant.LfiUrlPayins, constant.LfiUrlAvailablePayment)
	// operationPostRequestPayment for the ability to create specific payment with Merchant Order ID
	operationPostRequestPayment := fmt.Sprintf("%v /%v/%v", constant.LfiMethodPost, constant.LfiUrlPayins, id)
	// operationGetRequestPayment for the ability to get payment (but only with Merchant Order ID) details/status
	operationGetRequestPayment := fmt.Sprintf("%v /%v/%v", constant.LfiMethodGet, constant.LfiUrlPayins, id)

	// request token payload
	requestData := dto.TokenRequest{
		MerchantCode: cfg.MerchantCodeLfi,
		Secret:       cfg.SecretLfi,
		ValiditySecs: int(oneHour.Seconds()),
		Operations: []string{
			operationPostAvailablePayment,
			operationPostRequestPayment,
			operationGetRequestPayment,
		},
	}

	// transform to payload json
	payloadJson, _ := json.Marshal(requestData)

	// http-request
	r, err := http.NewRequest(http.MethodPost, cfg.TokenAuthLfiUrl, bytes.NewBuffer(payloadJson))
	if err != nil {
		return "", err
	}

	// set header
	r.Header.Add("Content-Type", "application/json")
	r.Header["X-API-Version"] = []string{cfg.ApiVersion}
	r.Close = true

	response, err := lfi.httpClient.Do(r)
	if err != nil {
		return "", err
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
		slog.Infof("LFI %v [provider-request] [request-token] got error failed to read response %v", id, string(contents))
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("LFI %v [provider-request] [request-token] got error response status not ok %v", id, string(contents))
		return "", errors.New("status not ok")
	}

	// convert response to struct
	var schema dto.TokenRequestResponse
	err = json.Unmarshal(contents, &schema)
	if err != nil {
		slog.Infof("LFI %v [provider-request] [request-token] got error failed unmarshal response %v", id, string(contents))
		return "", err
	}

	return schema.Token, nil
}

func (lfi *lfi) RequestVirtualAccount(payload dto.PayRequest) (string, error) {
	var schema dto.VaRequestResponse
	cfg := lfi.config

	token, err := lfi.RequestAuthToken(payload.Id)
	if err != nil {
		slog.Infof("LFI %v [provider-request] [payment-request] got error failed request token", payload.Id)
		return "", err
	}

	// bearer token
	bearer := "Bearer " + token

	// payment request payload
	requestData := dto.VaRequest{
		PaymentRequested: dto.PaymentRequestedWrapper{
			Money: dto.MoneyWrapper{
				Amount:       payload.Amount,
				CurrencyCode: constant.IDRCurrencyCode,
			},
		},
		PaymentMethod: dto.PaymentMethodWrapper{
			PaymentMethodCode: payload.BankCode,
			Account: dto.AccountWrapper{
				AccountName: payload.AccountName,
			},
		},
		CallbackUrl: cfg.CallbackUrl,
		ReturnUrl:   cfg.ReturnUrl,
	}

	// transform to payload json
	payloadJson, _ := json.Marshal(requestData)
	slog.Infof("LFI %v [provider-request] [payment-request] va request payload: %v", payload.Id, string(payloadJson))

	// transaction URL
	transactionUrl, _ := url.JoinPath(cfg.PaymentRequestLfiUrl, url.PathEscape(payload.Id))

	// http-request
	r, err := http.NewRequest(http.MethodPost, transactionUrl, bytes.NewBuffer(payloadJson))
	if err != nil {
		return "", err
	}

	// set header
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", bearer)
	r.Header["X-API-Version"] = []string{cfg.ApiVersion}
	r.Close = true

	response, err := lfi.httpClient.Do(r)
	if err != nil {
		return "", err
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
		slog.Infof("LFI %v [provider-request] [payment-request] got error failed to read response %v", payload.Id, string(contents))
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("LFI %v [provider-request] [payment-reques] got error response status not ok %v", payload.Id, string(contents))
		return "", errors.New("status not ok")
	}

	// convert response to struct
	err = json.Unmarshal(contents, &schema)
	if err != nil {
		slog.Infof("LFI %v [provider-request] [payment-request] got error failed unmarshal response %v", payload.Id, string(contents))
		return "", err
	}

	// va number
	virtualAccount := schema.PaymentMethodResponse.Account.AccountNumber

	return virtualAccount, nil
}

func (lfi *lfi) SyncStatus(id string) (dto.SyncStatusResponse, error) {
	var schema dto.SyncStatusResponse
	cfg := lfi.config

	token, err := lfi.RequestAuthToken(id)
	if err != nil {
		slog.Infof("LFI %v [provider-request] [sync-status] got error failed request token", id)
		return schema, err
	}

	// bearer token
	bearer := "Bearer " + token

	// transaction url
	transactionUrl, _ := url.JoinPath(cfg.PaymentRequestLfiUrl, url.PathEscape(id))

	// http-request
	r, err := http.NewRequest(http.MethodGet, transactionUrl, nil)
	if err != nil {
		return schema, err
	}

	// set header
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", bearer)
	r.Header["X-API-Version"] = []string{cfg.ApiVersion}
	r.Close = true

	response, err := lfi.httpClient.Do(r)
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
		slog.Infof("LFI %v [provider-request] [sync-status] got error failed to read response %v", id, string(contents))
		return schema, err
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("LFI %v [provider-request] [sync-status] got error response status not ok %v", id, string(contents))
		return schema, errors.New("status not ok")
	}

	// convert response status to struct
	err = json.Unmarshal(contents, &schema)
	if err != nil {
		slog.Infof("LFI %v [provider-request] [sync-status] got error failed unmarshal response %v", id, string(contents))
		return schema, err
	}

	if schema.Process.Status == constant.PaymentStatusExpired {
		slog.Infof("LFI %v [provider-request] [sync-status] got payment status not succeed: %v", id, converter.ToString(schema))
		return schema, nil
	}

	if schema.Process.Status != constant.PaymentStatusSuccess {
		slog.Infof("LFI %v [provider-request] [sync-status] got payment status not succeed: %v", id, converter.ToString(schema))
		return schema, nil
	}

	slog.Infof("LFI %v [provider-request] [sync-status] response data: %v", id, converter.ToString(schema))

	return schema, nil
}
