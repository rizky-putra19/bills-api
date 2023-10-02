package adapter

import (
	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/adapter/lfi"
	"gitlab.com/lokalpay-dev/digital-goods/internal/adapter/prismalink"
	"gitlab.com/lokalpay-dev/digital-goods/internal/adapter/tokopedia"
)

type Adapter struct {
	LFI        internal.LFIProvider
	Tokopedia  internal.TokopediaProvider
	Prismalink internal.PrismalinkProvider
}

func New(cfg config.Provider) *Adapter {
	return &Adapter{
		LFI:        lfi.New(cfg),
		Tokopedia:  tokopedia.New(cfg),
		Prismalink: prismalink.New(cfg),
	}
}
