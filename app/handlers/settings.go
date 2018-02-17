package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/tasks"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
)

// ChangeUserEmail register the intent of changing user e-mail
func ChangeUserEmail() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.ChangeUserEmail)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tenants.SaveVerificationKey(input.Model.VerificationKey, 24*time.Hour, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.SendChangeEmailConfirmation(input.Model))

		return c.Ok(web.Map{})
	}
}

// VerifyChangeEmailKey checks if key is correct and update user's email
func VerifyChangeEmailKey() web.HandlerFunc {
	return func(c web.Context) error {
		result, err := validateKey(models.EmailVerificationKindChangeEmail, c)
		if result == nil {
			return err
		}

		if result.UserID != c.User().ID {
			return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}

		err = c.Services().Users.ChangeEmail(result.UserID, result.Email)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Tenants.SetKeyAsVerified(result.Key)
		if err != nil {
			return c.Failure(err)
		}
		return c.Redirect(http.StatusTemporaryRedirect, c.BaseURL()+"/settings")
	}
}

// UserSettings is the current user's profile settings page
func UserSettings() web.HandlerFunc {
	return func(c web.Context) error {
		settings, err := c.Services().Users.GetUserSettings()
		if err != nil {
			return err
		}
		return c.Page(web.Map{
			"settings": settings,
		})
	}
}

// UpdateUserSettings updates current user settings
func UpdateUserSettings() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.UpdateUserSettings)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.Update(c.User().ID, input.Model)
		if err != nil {
			return c.Failure(err)
		}

		err = c.Services().Users.UpdateSettings(input.Model.Settings)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ChangeUserRole changes given user role
func ChangeUserRole() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.ChangeUserRole)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Users.ChangeRole(input.Model.UserID, input.Model.Role)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
