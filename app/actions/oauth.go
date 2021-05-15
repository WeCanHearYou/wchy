package actions

import (
	"context"
	"strings"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateEditOAuthConfig is used to create/edit OAuth config
type CreateEditOAuthConfig struct {
	ID                int
	Logo              *models.ImageUpload `json:"logo"`
	Provider          string              `json:"provider"`
	Status            int                 `json:"status"`
	DisplayName       string              `json:"displayName"`
	ClientID          string              `json:"clientID"`
	ClientSecret      string              `json:"clientSecret"`
	AuthorizeURL      string              `json:"authorizeURL" format:"lower"`
	TokenURL          string              `json:"tokenURL" format:"lower"`
	Scope             string              `json:"scope"`
	ProfileURL        string              `json:"profileURL" format:"lower"`
	JSONUserIDPath    string              `json:"jsonUserIDPath"`
	JSONUserNamePath  string              `json:"jsonUserNamePath"`
	JSONUserEmailPath string              `json:"jsonUserEmailPath"`
}

func NewCreateEditOAuthConfig() *CreateEditOAuthConfig {
	return &CreateEditOAuthConfig{
		Logo: &models.ImageUpload{},
	}
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateEditOAuthConfig) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateEditOAuthConfig) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if action.Provider != "" {
		getConfig := &query.GetCustomOAuthConfigByProvider{Provider: action.Provider}
		err := bus.Dispatch(ctx, getConfig)
		if err != nil {
			return validate.Error(err)
		}

		action.ID = getConfig.Result.ID
		action.Logo.BlobKey = getConfig.Result.LogoBlobKey
		if action.ClientSecret == "" {
			action.ClientSecret = getConfig.Result.ClientSecret
		}
	} else {
		action.Provider = "_" + strings.ToLower(rand.String(10))
	}

	messages, err := validate.ImageUpload(action.Logo, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    24,
		MinWidth:     24,
		ExactRatio:   true,
		MaxKilobytes: 50,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("logo", messages...)

	if action.Status != enum.OAuthConfigEnabled &&
		action.Status != enum.OAuthConfigDisabled {
		result.AddFieldFailure("status", "Invalid status.")
	}

	if action.DisplayName == "" {
		result.AddFieldFailure("displayName", "Display Name is required.")
	} else if len(action.DisplayName) > 50 {
		result.AddFieldFailure("displayName", "Display Name must have less than 50 characters.")
	}

	if action.ClientID == "" {
		result.AddFieldFailure("clientID", "Client ID is required.")
	} else if len(action.ClientID) > 100 {
		result.AddFieldFailure("clientID", "Client ID must have less than 100 characters.")
	}

	if action.ClientSecret == "" {
		result.AddFieldFailure("clientSecret", "Client Secret is required.")
	} else if len(action.ClientSecret) > 500 {
		result.AddFieldFailure("clientSecret", "Client Secret must have less than 500 characters.")
	}

	if action.Scope == "" {
		result.AddFieldFailure("scope", "Scope is required.")
	} else if len(action.Scope) > 100 {
		result.AddFieldFailure("scope", "Scope must have less than 100 characters.")
	}

	if action.AuthorizeURL == "" {
		result.AddFieldFailure("authorizeURL", "Authorize URL is required.")
	} else if messages := validate.URL(action.AuthorizeURL); len(messages) > 0 {
		result.AddFieldFailure("authorizeURL", messages...)
	}

	if action.TokenURL == "" {
		result.AddFieldFailure("tokenURL", "Token URL is required.")
	} else if messages := validate.URL(action.TokenURL); len(messages) > 0 {
		result.AddFieldFailure("tokenURL", messages...)
	}

	if action.ProfileURL != "" {
		if messages := validate.URL(action.ProfileURL); len(messages) > 0 {
			result.AddFieldFailure("profileURL", messages...)
		}
	}

	if action.JSONUserIDPath == "" {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path is required.")
	} else if len(action.JSONUserIDPath) > 100 {
		result.AddFieldFailure("jsonUserIDPath", "JSON User ID Path must have less than 100 characters.")
	}

	if len(action.JSONUserNamePath) > 100 {
		result.AddFieldFailure("jsonUserNamePath", "JSON User Name Path must have less than 100 characters.")
	}

	if len(action.JSONUserEmailPath) > 100 {
		result.AddFieldFailure("jsonUserEmailPath", "JSON User Email Path must have less than 100 characters.")
	}

	return result
}
