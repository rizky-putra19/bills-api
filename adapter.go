package adapter

import (
	"gitlab.com/lokalpay-dev/digital-goods/config"
)

type Adapter struct {
}

func New(cfg config.Provider) *Adapter {
	return &Adapter{}
}
