package postgres

import (
	"context"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
)

func getWebhook(ctx context.Context, q *query.GetWebhook) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		webhook := &entity.Webhook{}
		if tenant != nil {
			err := trx.Get(webhook, `
			SELECT id, name, type, status, url, content, http_method, additional_http_headers 
			FROM webhooks 
			WHERE tenant_id = $1 AND id = $2`, tenant.ID, q.ID)
			if err != nil {
				return err
			}
		}

		q.Result = webhook
		return nil
	})
}

func listAllWebhooks(ctx context.Context, q *query.ListAllWebhooks) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		webhooks := []*entity.Webhook{}
		if tenant != nil {
			err := trx.Select(&webhooks, `
			SELECT id, name, type, status, url, content, http_method, additional_http_headers 
			FROM webhooks 
			WHERE tenant_id = $1 
			ORDER BY id`, tenant.ID)
			if err != nil {
				return err
			}
		}

		q.Result = webhooks
		return nil
	})
}

func listAllWebhooksByType(ctx context.Context, q *query.ListAllWebhooksByType) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		webhooks := []*entity.Webhook{}
		if tenant != nil {
			err := trx.Select(&webhooks, `
			SELECT id, name, type, status, url, content, http_method, additional_http_headers 
			FROM webhooks 
			WHERE tenant_id = $1 AND type = $2 
			ORDER BY id`, tenant.ID, q.Type)
			if err != nil {
				return err
			}
		}

		q.Result = webhooks
		return nil
	})
}

func listActiveWebhooksByType(ctx context.Context, q *query.ListActiveWebhooksByType) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		webhooks := []*entity.Webhook{}
		if tenant != nil {
			err := trx.Select(&webhooks, `
			SELECT id, name, type, status, url, content, http_method, additional_http_headers 
			FROM webhooks 
			WHERE tenant_id = $1 AND type = $2 AND status = 1 
			ORDER BY id`, tenant.ID, q.Type)
			if err != nil {
				return err
			}
		}

		q.Result = webhooks
		return nil
	})
}

func createEditWebhook(ctx context.Context, q *query.CreateEditWebhook) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var err error
		id := q.ID

		if q.ID == 0 {
			err = trx.Get(&id, `
			INSERT INTO webhooks (name, type, status, url, content, http_method, additional_http_headers, tenant_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING id`, q.Name, q.Type, q.Status, q.Url, q.Content, q.HttpMethod, q.AdditionalHttpHeaders, tenant.ID)
		} else {
			_, err = trx.Execute(`
			UPDATE webhooks 
			SET name = $3, type = $4, status = $5, url = $6, content = $7, http_method = $8, additional_http_headers = $9 
			WHERE tenant_id = $1 AND id = $2`, tenant.ID, q.ID, q.Name, q.Type, q.Status, q.Url, q.Content, q.HttpMethod, q.AdditionalHttpHeaders)
		}

		if err != nil {
			return err
		}
		q.Result = id
		return nil
	})
}

func deleteWebhook(ctx context.Context, q *query.DeleteWebhook) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant != nil {
			_, err := trx.Execute(`
			DELETE FROM webhooks 
			WHERE tenant_id = $1 AND id = $2`, tenant.ID, q.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func markWebhookAsFailed(ctx context.Context, q *query.MarkWebhookAsFailed) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant != nil {
			_, err := trx.Execute(`
			UPDATE webhooks 
			SET status = $3 
			WHERE tenant_id = $1 AND id = $2`, tenant.ID, q.ID, enum.WebhookFailed)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
