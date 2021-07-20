package query

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type GetWebhook struct {
	ID int

	Result *entity.Webhook
}

type ListAllWebhooks struct {
	Result []*entity.Webhook
}

type ListAllWebhooksByType struct {
	Type string

	Result []*entity.Webhook
}

type ListActiveWebhooksByType struct {
	Type enum.WebhookType

	Result []*entity.Webhook
}

type CreateEditWebhook struct {
	ID                    int
	Name                  string
	Type                  enum.WebhookType
	Status                enum.WebhookStatus
	Url                   string
	Content               string
	HttpMethod            string
	AdditionalHttpHeaders entity.HttpHeaders

	Result int
}

type DeleteWebhook struct {
	ID int
}

type MarkWebhookAsFailed struct {
	ID int
}
