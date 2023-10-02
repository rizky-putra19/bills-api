package tokopedia

import (
	"net/http"

	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal/constant"
)

type tokopedia struct {
	config     config.Provider
	httpClient *http.Client
}

func New(cfg config.Provider) *tokopedia {
	return &tokopedia{
		config: cfg,
		httpClient: &http.Client{
			Timeout: constant.NinetySecond,
		},
	}
}
