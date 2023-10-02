package internal

import "gitlab.com/lokalpay-dev/digital-goods/internal/dto"

type LFIProvider interface {
	RequestAuthToken(id string) (string, error)
	RequestVirtualAccount(payload dto.PayRequest) (string, error)
	SyncStatus(id string) (dto.SyncStatusResponse, error)
}

type TokopediaProvider interface {
}

type PrismalinkProvider interface {
	RequestBiller(payload dto.RequestPayloadPrismalink, reqType string) (dto.RequestBillerResponsePayload, error) // reqType can filled INQUIRY	or PAYMENT
}
