package webhook

import (
	"context"
	"fmt"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/tpl"
	"net/http"
	"strings"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Webhook"
}

func (s Service) Category() string {
	return "hook"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(triggerWebhook)
	bus.AddHandler(triggerWebhooksByType)
	bus.AddHandler(previewWebhook)
	bus.AddHandler(getWebhookProps)
}

func triggerWebhook(ctx context.Context, c *cmd.TriggerWebhook) error {
	webhook := &query.GetWebhook{ID: c.ID}
	err := bus.Dispatch(ctx, webhook)
	if err != nil {
		return err
	}

	c.Result, err = runTriggerWebhook(ctx, webhook.Result, dummyTriggerProps(ctx, webhook.Result.Type))
	if err != nil {
		return err
	}

	return nil
}

func triggerWebhooksByType(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
	webhooks := &query.ListActiveWebhooksByType{Type: c.Type}
	err := bus.Dispatch(ctx, webhooks)
	if err != nil {
		return err
	}

	for _, webhook := range webhooks.Result {
		_, err = runTriggerWebhook(ctx, webhook, c.Props)
		if err != nil {
			return err
		}
	}

	return nil
}

func runTriggerWebhook(ctx context.Context, webhook *entity.Webhook, props dto.Props) (*entity.WebhookTriggerResult, error) {
	result := &entity.WebhookTriggerResult{Webhook: webhook, Props: props}
	var err error

	fullName := fmt.Sprintf("%d-%s", webhook.ID, webhook.Name)
	result.Url, err = executeTemplate(fmt.Sprintf("%s-url", fullName), webhook.Url, props)
	if err != nil {
		return resultWithError(ctx, "Could not parse webhook URL template", err.Error(), result)
	}
	result.Content, err = executeTemplate(fmt.Sprintf("%s-content", fullName), webhook.Content, props)
	if err != nil {
		return resultWithError(ctx, "Could not parse webhook content template", err.Error(), result)
	}

	httpRequest := &cmd.HTTPRequest{
		URL:       result.Url,
		Body:      strings.NewReader(result.Content),
		Method:    webhook.HttpMethod,
		Headers:   webhook.AdditionalHttpHeaders,
		BasicAuth: nil,
	}
	err = bus.Dispatch(ctx, httpRequest)
	if err != nil {
		return resultWithError(ctx, "Could not execute webhook HTTP request", err.Error(), result)
	}
	result.StatusCode = httpRequest.ResponseStatusCode
	if result.StatusCode >= http.StatusBadRequest {
		fullResponse := fmt.Sprintf("%d %s:\n%s", result.StatusCode, http.StatusText(result.StatusCode), httpRequest.ResponseBody)
		return resultWithError(ctx, "Webhook HTTP request returned an error response code", fullResponse, result)
	}

	result.Success = true
	log.Infof(ctx, "Webhook #@{ID:yellow} @{Name:blue} finished with @{Code:magenta}", dto.Props{
		"ID":   webhook.ID,
		"Name": webhook.Name,
		"Code": result.StatusCode,
	})
	return result, nil
}

func previewWebhook(ctx context.Context, c *cmd.PreviewWebhook) error {
	c.Result = &entity.WebhookPreviewResult{}
	var err error
	props := dummyTriggerProps(ctx, c.Type)

	c.Result.Url.Value, err = executeTemplate("preview-url", c.Url, props)
	if err != nil {
		c.Result.Url.Message = "Could not parse webhook URL template"
		c.Result.Url.Error = err.Error()
		// Do not propagate error: it's a preview
	}
	c.Result.Content.Value, err = executeTemplate("preview-content", c.Content, props)
	if err != nil {
		c.Result.Content.Message = "Could not parse webhook content template"
		c.Result.Content.Error = err.Error()
		// Do not propagate error: it's a preview
	}

	return nil
}

func executeTemplate(name, text string, props dto.Props) (string, error) {
	tmpl, err := tpl.GetTextTemplate(name, text)
	if err != nil {
		return "", err
	}

	replacedText, err := tpl.Execute(tmpl, props)
	if err != nil {
		return "", err
	}

	return replacedText, nil
}

func resultWithError(ctx context.Context, message, error string, result *entity.WebhookTriggerResult) (*entity.WebhookTriggerResult, error) {
	result.Success = false
	result.Message = message
	result.Error = error

	log.Warnf(ctx, "@{Message} (ID: @{ID:yellow}, Name: @{Name:blue}): @{Error:red}", dto.Props{
		"Message": message,
		"ID":      result.Webhook.ID,
		"Name":    result.Webhook.Name,
		"Error":   error,
	})

	webhooks := &query.MarkWebhookAsFailed{ID: result.Webhook.ID}
	err := bus.Dispatch(ctx, webhooks)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getWebhookProps(ctx context.Context, c *cmd.GetWebhookProps) error {
	c.Result = dummyTriggerProps(ctx, c.Type)
	return nil
}
