package webhook

import (
	"context"
	"fmt"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/web"
	"time"
)

var dummyPost = &entity.Post{
	ID:     36,
	Number: 36,
	Title:  "Example dummy post title",
	Slug:   "example-dummy-post-title",
	Description: `*Those data are fake.*

This is an example of a post, but don't worry, nothing was created on your Fider instance:
It was just a test of a webhook trigger.

You can use ` + "`" + `{{ markdown .post_description }}` + "`" + ` to parse __Markdown__ to **HTML** in post!`,
	CreatedAt: time.Date(2021, time.May, 7, 18, 42, 27, 0, time.UTC),
	User: &entity.User{
		ID:    7,
		Name:  "Fider",
		Email: "contact@fider.io",
		Role:  1,
	},
	HasVoted:      true,
	VotesCount:    7,
	CommentsCount: 3,
	Status:        enum.PostStarted,
	Response: &entity.PostResponse{
		Text:        "This is a response, still in *Markdown*.",
		RespondedAt: time.Date(2021, time.July, 9, 15, 29, 57, 0, time.UTC),
		User:        nil,
	},
	Tags: []string{"tag1", "tag2"},
}

func dummyTriggerProps(c context.Context, webhookType enum.WebhookType) dto.Props {
	props := dto.Props{}
	author := c.Value(app.UserCtxKey).(*entity.User)
	tenant := c.Value(app.TenantCtxKey).(*entity.Tenant)
	baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)
	dummyPost.User.AvatarURL = logoURL
	dummyPost.Response.User = author
	DescribeUser(props, author, "author")
	DescribeTenant(props, tenant, "tenant", baseURL, logoURL)
	switch webhookType {
	case enum.WebhookNewPost:
		DescribePost(props, dummyPost, "post", baseURL, false, false)
	case enum.WebhookNewComment:
		DescribePost(props, dummyPost, "post", baseURL, true, true)
		props["comment"] = "An example **comment** on a post."
	case enum.WebhookChangeStatus:
		DescribePost(props, dummyPost, "post", baseURL, true, true)
		props["post_old_status"] = enum.PostOpen
	case enum.WebhookDeletePost:
		DescribePost(props, dummyPost, "post", baseURL, true, true)
		props["post_status"] = enum.PostDeleted
		props["post_response_text"] = "The reason _why_ this post was deleted."
	}
	return props
}

func DescribeUser(props dto.Props, user *entity.User, keyPrefix string) {
	if user == nil {
		return
	}
	role, _ := user.Role.MarshalText()
	props.Append(dto.Props{
		keyPrefix + "_id":     user.ID,
		keyPrefix + "_name":   user.Name,
		keyPrefix + "_email":  user.Email,
		keyPrefix + "_role":   string(role),
		keyPrefix + "_avatar": user.AvatarURL,
	})
}

func DescribePost(props dto.Props, post *entity.Post, keyPrefix, baseURL string, complete, author bool) {
	if post == nil {
		return
	}
	props.Append(dto.Props{
		keyPrefix + "_id":          post.ID,
		keyPrefix + "_number":      post.Number,
		keyPrefix + "_title":       post.Title,
		keyPrefix + "_slug":        post.Slug,
		keyPrefix + "_description": post.Description,
		keyPrefix + "_created_at":  post.CreatedAt,
		keyPrefix + "_url":         fmt.Sprintf("%s/posts/%d/%s", baseURL, post.Number, post.Slug),
	})
	if author {
		DescribeUser(props, post.User, keyPrefix+"_author")
	}
	if complete {
		postResponse := post.Response
		props.Append(dto.Props{
			keyPrefix + "_votes":    post.VotesCount,
			keyPrefix + "_comments": post.CommentsCount,
			keyPrefix + "_status":   post.Status.Name(),
			keyPrefix + "_tags":     post.Tags,
			keyPrefix + "_response": postResponse != nil,
		})
		if postResponse != nil {
			keyPrefix := keyPrefix + "_response"
			props.Append(dto.Props{
				keyPrefix + "_text":         postResponse.Text,
				keyPrefix + "_responded_at": postResponse.RespondedAt,
			})
			DescribeUser(props, postResponse.User, keyPrefix+"_author")
			originalPost := postResponse.Original
			if post.Status == enum.PostDuplicate && originalPost != nil {
				keyPrefix := keyPrefix + "_original"
				props.Append(dto.Props{
					keyPrefix + "_number": originalPost.Number,
					keyPrefix + "_title":  originalPost.Title,
					keyPrefix + "_slug":   originalPost.Slug,
					keyPrefix + "_status": originalPost.Status.Name(),
					keyPrefix + "_url":    fmt.Sprintf("%s/posts/%d/%s", baseURL, originalPost.Number, originalPost.Slug),
				})
			}
		}
	}
}

func DescribeTenant(props dto.Props, tenant *entity.Tenant, keyPrefix, baseURL, logoURL string) {
	if tenant == nil {
		return
	}
	props.Append(dto.Props{
		keyPrefix + "_id":              tenant.ID,
		keyPrefix + "_name":            tenant.Name,
		keyPrefix + "_subdomain":       tenant.Subdomain,
		keyPrefix + "_invitation":      tenant.Invitation,
		keyPrefix + "_welcome_message": tenant.WelcomeMessage,
		keyPrefix + "_status":          TenantStatus(tenant.Status),
		keyPrefix + "_locale":          tenant.Locale,
		keyPrefix + "_url":             baseURL,
		keyPrefix + "_logo":            logoURL,
	})
}

func TenantStatus(status int) string {
	switch status {
	case enum.TenantActive:
		return "active"
	case enum.TenantPending:
		return "pending"
	case enum.TenantLocked:
		return "locked"
	case enum.TenantDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}
