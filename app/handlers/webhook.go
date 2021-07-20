package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"strconv"
)

// ManageWebhooks is the page used by administrators to configure webhooks
func ManageWebhooks() web.HandlerFunc {
	return func(c *web.Context) error {
		allWebhooks := &query.ListAllWebhooks{}
		if err := bus.Dispatch(c, allWebhooks); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Manage Webhooks · Site Settings",
			ChunkName: "ManageWebhooks.page",
			Data: web.Map{
				"webhooks": allWebhooks.Result,
			},
		})
	}
}

func CreateWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		action := &actions.CreateEditWebhook{}
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		createWebhook := &query.CreateEditWebhook{
			ID:                    0,
			Name:                  action.Name,
			Type:                  action.Type,
			Status:                action.Status,
			Url:                   action.Url,
			Content:               action.Content,
			HttpMethod:            action.HttpMethod,
			AdditionalHttpHeaders: action.AdditionalHttpHeaders,
		}
		if err := bus.Dispatch(c, createWebhook); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{"id": createWebhook.Result})
	}
}

func UpdateWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)
		if err != nil {
			return c.Failure(err)
		}

		action := &actions.CreateEditWebhook{}
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updateWebhook := &query.CreateEditWebhook{
			ID:                    id,
			Name:                  action.Name,
			Type:                  action.Type,
			Status:                action.Status,
			Url:                   action.Url,
			Content:               action.Content,
			HttpMethod:            action.HttpMethod,
			AdditionalHttpHeaders: action.AdditionalHttpHeaders,
		}
		if action.Status == enum.WebhookFailed {
			updateWebhook.Status = enum.WebhookDisabled
		}
		if err := bus.Dispatch(c, updateWebhook); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

func DeleteWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)
		if err != nil {
			return c.Failure(err)
		}

		deleteWebhook := &query.DeleteWebhook{ID: id}
		if err = bus.Dispatch(c, deleteWebhook); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

func TriggerWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)
		if err != nil {
			return c.Failure(err)
		}

		triggerWebhook := &cmd.TriggerWebhook{ID: id}
		if err = bus.Dispatch(c, triggerWebhook); err != nil {
			return c.Failure(err)
		}

		return c.Ok(triggerWebhook.Result)
	}
}

func PreviewWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		action := &actions.PreviewWebhook{}
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		previewWebhook := &cmd.PreviewWebhook{
			Type:    action.Type,
			Url:     action.Url,
			Content: action.Content,
		}
		if err := bus.Dispatch(c, previewWebhook); err != nil {
			return c.Failure(err)
		}

		return c.Ok(previewWebhook.Result)
	}
}

func GetWebhookProps() web.HandlerFunc {
	return func(c *web.Context) error {
		rawType := c.Param("type")
		var webhookType enum.WebhookType
		err := webhookType.UnmarshalText([]byte(rawType))
		if err != nil {
			return c.Failure(err)
		}

		webhookProps := &cmd.GetWebhookProps{Type: webhookType}
		if err := bus.Dispatch(c, webhookProps); err != nil {
			return c.Failure(err)
		}

		return c.Ok(webhookProps.Result)
	}
}
