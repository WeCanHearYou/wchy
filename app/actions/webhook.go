package actions

import (
	"context"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/validate"
)

type CreateEditWebhook struct {
	Name                  string             `json:"name"`
	Type                  enum.WebhookType   `json:"type"`
	Status                enum.WebhookStatus `json:"status"`
	Url                   string             `json:"url"`
	Content               string             `json:"content"`
	HttpMethod            string             `json:"http_method"`
	AdditionalHttpHeaders entity.HttpHeaders `json:"additional_http_headers"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateEditWebhook) IsAuthorized(_ context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateEditWebhook) Validate(ctx context.Context, _ *entity.User) *validate.Result {
	result := validate.Success()

	if action.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(action.Name) > 100 {
		result.AddFieldFailure("name", "Name must have less than 100 characters.")
	}

	if action.Type == 0 {
		result.AddFieldFailure("type", "Type is required.")
	} else if action.Type != enum.WebhookNewPost &&
		action.Type != enum.WebhookNewComment &&
		action.Type != enum.WebhookChangeStatus &&
		action.Type != enum.WebhookDeletePost {
		result.AddFieldFailure("type", "Type must be valid.")
	}

	if action.Status == 0 {
		result.AddFieldFailure("status", "Status is required.")
	}

	if action.Url == "" {
		result.AddFieldFailure("url", "URL is required.")
	} else if messages := validate.URL(ctx, action.Url); len(messages) > 0 {
		result.AddFieldFailure("url", messages...)
	}

	if action.HttpMethod == "" {
		result.AddFieldFailure("http_method", "HTTP Method is required.")
	} else if len(action.HttpMethod) > 50 {
		result.AddFieldFailure("http_method", "HTTP Method must have less than 50 characters.")
	}

	if len(action.Content) > 10_000 {
		result.AddFieldFailure("content", "Content must have less than 10 000 characters.")
	}

	for header, value := range action.AdditionalHttpHeaders {
		if header == "" {
			result.AddFieldFailure("header-"+header, "HTTP Header Name is required.")
		} else if len(header) > 200 {
			result.AddFieldFailure("header-"+header, "HTTP Header Name must have less than 200 characters.")
		}

		if value == "" {
			result.AddFieldFailure("value-"+header, "HTTP Header Value is required.")
		} else if len(value) > 1_000 {
			result.AddFieldFailure("value-"+header, "HTTP Header Value must have less than 1 000 characters.")
		}
	}

	return result
}

type PreviewWebhook struct {
	Type    enum.WebhookType `json:"type"`
	Url     string           `json:"url"`
	Content string           `json:"content"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *PreviewWebhook) IsAuthorized(_ context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *PreviewWebhook) Validate(context.Context, *entity.User) *validate.Result {
	result := validate.Success()

	if action.Type == 0 {
		result.AddFieldFailure("type", "Type is required.")
	} else if action.Type != enum.WebhookNewPost &&
		action.Type != enum.WebhookNewComment &&
		action.Type != enum.WebhookChangeStatus &&
		action.Type != enum.WebhookDeletePost {
		result.AddFieldFailure("type", "Type must be valid.")
	}

	return result
}
