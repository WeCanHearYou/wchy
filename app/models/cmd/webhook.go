package cmd

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type TriggerWebhook struct {
	ID int

	Result *entity.WebhookTriggerResult
}

type TriggerWebhooksByType struct {
	Type  enum.WebhookType
	Props dto.Props
}

type PreviewWebhook struct {
	Type    enum.WebhookType
	Url     string
	Content string

	Result *entity.WebhookPreviewResult
}

type GetWebhookProps struct {
	Type enum.WebhookType

	Result dto.Props
}
