package tasks_test

import (
	"context"
	"github.com/getfider/fider/app/services/webhook"
	"html/template"
	"testing"
	"time"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/services/email/emailmock"
	"github.com/getfider/fider/app/tasks"
)

func TestSendSignUpEmailTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendSignUpEmail(&actions.CreateTenant{
		VerificationKey: "1234",
	}, "http://domain.com")

	err := worker.
		AsUser(mock.JonSnow).
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("signup_email")
	Expect(emailmock.MessageHistory[0].Tenant).IsNil()
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Fider")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Props: dto.Props{
			"link": "<a href='http://domain.com/signup/verify?k=1234'>http://domain.com/signup/verify?k=1234</a>",
		},
	})
}

func TestSendSignInEmailTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendSignInEmail("jon@got.com", "9876")

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("signin_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Address: "jon@got.com",
		Props: dto.Props{
			"siteName": mock.DemoTenant.Name,
			"link":     "<a href='http://domain.com/signin/verify?k=9876'>http://domain.com/signin/verify?k=9876</a>",
		},
	})
}

func TestSendChangeEmailConfirmationTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	worker := mock.NewWorker()
	task := tasks.SendChangeEmailConfirmation(&actions.ChangeUserEmail{
		Email:           "newemail@domain.com",
		VerificationKey: "13579",
		Requestor:       mock.JonSnow,
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_emailaddress_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"logo": "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals(mock.DemoTenant.Name)
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "newemail@domain.com",
		Props: dto.Props{
			"name":     "Jon Snow",
			"oldEmail": "jon.snow@got.com",
			"newEmail": "newemail@domain.com",
			"link":     "<a href='http://domain.com/change-email/verify?k=13579'>http://domain.com/change-email/verify?k=13579</a>",
		},
	})
}

func TestNotifyAboutNewPostTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	var triggerWebhooksByType *cmd.TriggerWebhooksByType
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
		triggerWebhooksByType = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
	}
	task := tasks.NotifyAboutNewPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_post")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":    "Add support for TypeScript",
		"postLink": "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName": "Demonstration",
		"userName": "Jon Snow",
		"content":  template.HTML("<p>TypeScript is great, please add support for it</p>"),
		"view":     "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":   "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"logo":     "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("New post: **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooksByType).IsNotNil()
	Expect(triggerWebhooksByType.Type).Equals(enum.WebhookNewPost)
	role, _ := mock.JonSnow.Role.MarshalText()
	Expect(triggerWebhooksByType.Props).ContainsProps(dto.Props{
		"post_id":          post.ID,
		"post_number":      post.Number,
		"post_title":       post.Title,
		"post_slug":        post.Slug,
		"post_description": post.Description,
		"post_url":         "http://domain.com/posts/1/add-support-for-typescript",
		"author_id":        mock.JonSnow.ID,
		"author_name":      mock.JonSnow.Name,
		"author_email":     mock.JonSnow.Email,
		"author_role":      string(role),
		"tenant_id":        mock.DemoTenant.ID,
		"tenant_name":      mock.DemoTenant.Name,
		"tenant_subdomain": mock.DemoTenant.Subdomain,
		"tenant_status":    webhook.TenantStatus(mock.DemoTenant.Status),
		"tenant_url":       "http://domain.com",
	})
}

func TestNotifyAboutNewCommentTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.JonSnow,
		}
		return nil
	})

	var triggerWebhooksByType *cmd.TriggerWebhooksByType
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
		triggerWebhooksByType = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.JonSnow,
	}
	task := tasks.NotifyAboutNewComment(post, "I agree")

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("new_comment")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":       "Add support for TypeScript",
		"postLink":    "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":    "Demonstration",
		"userName":    "Arya Stark",
		"content":     template.HTML("<p>I agree</p>"),
		"view":        "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":      "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe": "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Arya Stark")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Jon Snow",
		Address: "jon.snow@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Arya Stark** left a comment on **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.JonSnow)

	Expect(triggerWebhooksByType).IsNotNil()
	Expect(triggerWebhooksByType.Type).Equals(enum.WebhookNewComment)
	roleArya, _ := mock.AryaStark.Role.MarshalText()
	roleJon, _ := mock.JonSnow.Role.MarshalText()
	Expect(triggerWebhooksByType.Props).ContainsProps(dto.Props{
		"comment":           "I agree",
		"post_id":           post.ID,
		"post_number":       post.Number,
		"post_title":        post.Title,
		"post_slug":         post.Slug,
		"post_description":  post.Description,
		"post_url":          "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":    mock.JonSnow.ID,
		"post_author_name":  mock.JonSnow.Name,
		"post_author_email": mock.JonSnow.Email,
		"post_author_role":  string(roleJon),
		"author_id":         mock.AryaStark.ID,
		"author_name":       mock.AryaStark.Name,
		"author_email":      mock.AryaStark.Email,
		"author_role":       string(roleArya),
		"tenant_id":         mock.DemoTenant.ID,
		"tenant_name":       mock.DemoTenant.Name,
		"tenant_subdomain":  mock.DemoTenant.Subdomain,
		"tenant_status":     webhook.TenantStatus(mock.DemoTenant.Status),
		"tenant_url":        "http://domain.com",
	})
}

func TestNotifyAboutStatusChangeTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	var triggerWebhooksByType *cmd.TriggerWebhooksByType
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
		triggerWebhooksByType = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.AryaStark,
		Status:      enum.PostPlanned,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			Text:        "Planned for next release.",
			User:        mock.JonSnow,
		},
	}

	task := tasks.NotifyAboutStatusChange(post, enum.PostOpen)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_status")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":       "Add support for TypeScript",
		"postLink":    "<a href='http://domain.com/posts/1/add-support-for-typescript'>#1</a>",
		"siteName":    "Demonstration",
		"content":     template.HTML("<p>Planned for next release.</p>"),
		"duplicate":   "",
		"status":      "Planned",
		"view":        "<a href='http://domain.com/posts/1/add-support-for-typescript'>view it on your browser</a>",
		"change":      "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe": "<a href='http://domain.com/posts/1/add-support-for-typescript'>unsubscribe from it</a>",
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/1/add-support-for-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **Add support for TypeScript** to **planned**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooksByType).IsNotNil()
	Expect(triggerWebhooksByType.Type).Equals(enum.WebhookChangeStatus)
	roleJon, _ := mock.JonSnow.Role.MarshalText()
	roleArya, _ := mock.AryaStark.Role.MarshalText()
	Expect(triggerWebhooksByType.Props).ContainsProps(dto.Props{
		"post_old_status":            enum.PostOpen,
		"post_id":                    post.ID,
		"post_number":                post.Number,
		"post_title":                 post.Title,
		"post_slug":                  post.Slug,
		"post_description":           post.Description,
		"post_status":                post.Status.Name(),
		"post_url":                   "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":             mock.AryaStark.ID,
		"post_author_name":           mock.AryaStark.Name,
		"post_author_email":          mock.AryaStark.Email,
		"post_author_role":           string(roleArya),
		"post_response":              true,
		"post_response_text":         post.Response.Text,
		"post_response_responded_at": post.Response.RespondedAt,
		"post_response_author_id":    mock.JonSnow.ID,
		"post_response_author_name":  mock.JonSnow.Name,
		"post_response_author_email": mock.JonSnow.Email,
		"post_response_author_role":  string(roleJon),
		"author_id":                  mock.JonSnow.ID,
		"author_name":                mock.JonSnow.Name,
		"author_email":               mock.JonSnow.Email,
		"author_role":                string(roleJon),
		"tenant_id":                  mock.DemoTenant.ID,
		"tenant_name":                mock.DemoTenant.Name,
		"tenant_subdomain":           mock.DemoTenant.Subdomain,
		"tenant_status":              webhook.TenantStatus(mock.DemoTenant.Status),
		"tenant_url":                 "http://domain.com",
	})
}

func TestNotifyAboutDeletePostTask(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	var triggerWebhooksByType *cmd.TriggerWebhooksByType
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
		triggerWebhooksByType = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:          1,
		Number:      1,
		Title:       "Add support for TypeScript",
		Slug:        "add-support-for-typescript",
		Description: "TypeScript is great, please add support for it",
		User:        mock.AryaStark,
		Status:      enum.PostDeleted,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			Text:        "Invalid post!",
			User:        mock.JonSnow,
		},
	}

	task := tasks.NotifyAboutDeletedPost(post)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("delete_post")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":    "Add support for TypeScript",
		"siteName": "Demonstration",
		"content":  template.HTML("<p>Invalid post!</p>"),
		"change":   "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"logo":     "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("")
	Expect(addNewNotification.Title).Equals("**Jon Snow** deleted **Add support for TypeScript**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooksByType).IsNotNil()
	Expect(triggerWebhooksByType.Type).Equals(enum.WebhookDeletePost)
	roleJon, _ := mock.JonSnow.Role.MarshalText()
	roleArya, _ := mock.AryaStark.Role.MarshalText()
	Expect(triggerWebhooksByType.Props).ContainsProps(dto.Props{
		"post_id":                    post.ID,
		"post_number":                post.Number,
		"post_title":                 post.Title,
		"post_slug":                  post.Slug,
		"post_description":           post.Description,
		"post_status":                post.Status.Name(),
		"post_url":                   "http://domain.com/posts/1/add-support-for-typescript",
		"post_author_id":             mock.AryaStark.ID,
		"post_author_name":           mock.AryaStark.Name,
		"post_author_email":          mock.AryaStark.Email,
		"post_author_role":           string(roleArya),
		"post_response":              true,
		"post_response_text":         post.Response.Text,
		"post_response_responded_at": post.Response.RespondedAt,
		"post_response_author_id":    mock.JonSnow.ID,
		"post_response_author_name":  mock.JonSnow.Name,
		"post_response_author_email": mock.JonSnow.Email,
		"post_response_author_role":  string(roleJon),
		"author_id":                  mock.JonSnow.ID,
		"author_name":                mock.JonSnow.Name,
		"author_email":               mock.JonSnow.Email,
		"author_role":                string(roleJon),
		"tenant_id":                  mock.DemoTenant.ID,
		"tenant_name":                mock.DemoTenant.Name,
		"tenant_subdomain":           mock.DemoTenant.Subdomain,
		"tenant_status":              webhook.TenantStatus(mock.DemoTenant.Status),
		"tenant_url":                 "http://domain.com",
	})
}

func TestNotifyAboutStatusChangeTask_Duplicate(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	var addNewNotification *cmd.AddNewNotification
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewNotification) error {
		addNewNotification = c
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetActiveSubscribers) error {
		q.Result = []*entity.User{
			mock.AryaStark,
		}
		return nil
	})

	var triggerWebhooksByType *cmd.TriggerWebhooksByType
	bus.AddHandler(func(ctx context.Context, c *cmd.TriggerWebhooksByType) error {
		triggerWebhooksByType = c
		return nil
	})

	worker := mock.NewWorker()
	post := &entity.Post{
		ID:     2,
		Number: 2,
		Title:  "I need TypeScript",
		Slug:   "i-need-typescript",
		User:   mock.AryaStark,
		Status: enum.PostDuplicate,
		Response: &entity.PostResponse{
			RespondedAt: time.Now(),
			User:        mock.JonSnow,
			Original: &entity.OriginalPost{
				Number: 1,
				Title:  "Add support for TypeScript",
				Slug:   "add-support-for-typescript",
				Status: enum.PostOpen,
			},
		},
	}

	task := tasks.NotifyAboutStatusChange(post, enum.PostOpen)

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("change_status")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"title":       "I need TypeScript",
		"postLink":    "<a href='http://domain.com/posts/2/i-need-typescript'>#2</a>",
		"siteName":    "Demonstration",
		"content":     template.HTML(""),
		"duplicate":   "<a href='http://domain.com/posts/1/add-support-for-typescript'>Add support for TypeScript</a>",
		"status":      "Duplicate",
		"view":        "<a href='http://domain.com/posts/2/i-need-typescript'>view it on your browser</a>",
		"change":      "<a href='http://domain.com/settings'>change your notification preferences</a>",
		"unsubscribe": "<a href='http://domain.com/posts/2/i-need-typescript'>unsubscribe from it</a>",
		"logo":        "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(1)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Name:    "Arya Stark",
		Address: "arya.stark@got.com",
		Props:   dto.Props{},
	})

	Expect(addNewNotification).IsNotNil()
	Expect(addNewNotification.PostID).Equals(post.ID)
	Expect(addNewNotification.Link).Equals("/posts/2/i-need-typescript")
	Expect(addNewNotification.Title).Equals("**Jon Snow** changed status of **I need TypeScript** to **duplicate**")
	Expect(addNewNotification.User).Equals(mock.AryaStark)

	Expect(triggerWebhooksByType).IsNotNil()
	Expect(triggerWebhooksByType.Type).Equals(enum.WebhookChangeStatus)
	roleJon, _ := mock.JonSnow.Role.MarshalText()
	roleArya, _ := mock.AryaStark.Role.MarshalText()
	Expect(triggerWebhooksByType.Props).ContainsProps(dto.Props{
		"post_old_status":               enum.PostOpen,
		"post_id":                       post.ID,
		"post_number":                   post.Number,
		"post_title":                    post.Title,
		"post_slug":                     post.Slug,
		"post_description":              post.Description,
		"post_status":                   post.Status.Name(),
		"post_url":                      "http://domain.com/posts/2/i-need-typescript",
		"post_author_id":                mock.AryaStark.ID,
		"post_author_name":              mock.AryaStark.Name,
		"post_author_email":             mock.AryaStark.Email,
		"post_author_role":              string(roleArya),
		"post_response":                 true,
		"post_response_text":            post.Response.Text,
		"post_response_responded_at":    post.Response.RespondedAt,
		"post_response_author_id":       mock.JonSnow.ID,
		"post_response_author_name":     mock.JonSnow.Name,
		"post_response_author_email":    mock.JonSnow.Email,
		"post_response_author_role":     string(roleJon),
		"post_response_original_number": post.Response.Original.Number,
		"post_response_original_title":  post.Response.Original.Title,
		"post_response_original_slug":   post.Response.Original.Slug,
		"post_response_original_status": post.Response.Original.Status.Name(),
		"post_response_original_url":    "http://domain.com/posts/1/add-support-for-typescript",
		"author_id":                     mock.JonSnow.ID,
		"author_name":                   mock.JonSnow.Name,
		"author_email":                  mock.JonSnow.Email,
		"author_role":                   string(roleJon),
		"tenant_id":                     mock.DemoTenant.ID,
		"tenant_name":                   mock.DemoTenant.Name,
		"tenant_subdomain":              mock.DemoTenant.Subdomain,
		"tenant_status":                 webhook.TenantStatus(mock.DemoTenant.Status),
		"tenant_url":                    "http://domain.com",
	})
}

func TestSendInvites(t *testing.T) {
	RegisterT(t)
	bus.Init(emailmock.Service{})

	savedKeys := make([]*cmd.SaveVerificationKey, 0)
	bus.AddHandler(func(ctx context.Context, c *cmd.SaveVerificationKey) error {
		savedKeys = append(savedKeys, c)
		return nil
	})

	worker := mock.NewWorker()
	task := tasks.SendInvites("My Subject", "Click here: %invite%", []*actions.UserInvitation{
		{Email: "user1@domain.com", VerificationKey: "1234"},
		{Email: "user2@domain.com", VerificationKey: "5678"},
	})

	err := worker.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithBaseURL("http://domain.com").
		Execute(task)

	Expect(err).IsNil()
	Expect(emailmock.MessageHistory).HasLen(1)
	Expect(emailmock.MessageHistory[0].TemplateName).Equals("invite_email")
	Expect(emailmock.MessageHistory[0].Tenant).Equals(mock.DemoTenant)
	Expect(emailmock.MessageHistory[0].Props).Equals(dto.Props{
		"subject": "My Subject",
		"logo":    "https://getfider.com/images/logo-100x100.png",
	})
	Expect(emailmock.MessageHistory[0].From).Equals("Jon Snow")
	Expect(emailmock.MessageHistory[0].To).HasLen(2)
	Expect(emailmock.MessageHistory[0].To[0]).Equals(dto.Recipient{
		Address: "user1@domain.com",
		Props: dto.Props{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=1234">http://domain.com/invite/verify?k=1234</a></p>`),
		},
	})
	Expect(emailmock.MessageHistory[0].To[1]).Equals(dto.Recipient{
		Address: "user2@domain.com",
		Props: dto.Props{
			"message": template.HTML(`<p>Click here: <a href="http://domain.com/invite/verify?k=5678">http://domain.com/invite/verify?k=5678</a></p>`),
		},
	})

	Expect(savedKeys).HasLen(2)
	Expect(savedKeys[0].Key).Equals("1234")
	Expect(savedKeys[0].Request.GetEmail()).Equals("user1@domain.com")
	Expect(savedKeys[1].Key).Equals("5678")
	Expect(savedKeys[1].Request.GetEmail()).Equals("user2@domain.com")
}
