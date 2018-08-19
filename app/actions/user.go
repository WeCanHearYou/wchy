package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

//ChangeUserRole is the input model change role of an user
type ChangeUserRole struct {
	Model *models.ChangeUserRole
}

// Initialize the model
func (input *ChangeUserRole) Initialize() interface{} {
	input.Model = new(models.ChangeUserRole)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *ChangeUserRole) IsAuthorized(user *models.User, services *app.Services) bool {
	if user == nil {
		return false
	}
	return user.IsAdministrator() && user.ID != input.Model.UserID
}

// Validate is current model is valid
func (input *ChangeUserRole) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()
	if input.Model.Role < models.RoleVisitor || input.Model.Role > models.RoleAdministrator {
		return validate.Error(app.ErrNotFound)
	}

	if user.ID == input.Model.UserID {
		result.AddFieldFailure("userId", "It is not allowed to change your own Role.")
	}

	target, err := services.Users.GetByID(input.Model.UserID)
	if err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			result.AddFieldFailure("userId", "User not found.")
		} else {
			return validate.Error(err)
		}
	} else if target.Tenant.ID != user.Tenant.ID {
		result.AddFieldFailure("userId", "User not found.")
	}
	return result
}

//ChangeUserEmail is the action used to change current user's email
type ChangeUserEmail struct {
	Model *models.ChangeUserEmail
}

// Initialize the model
func (input *ChangeUserEmail) Initialize() interface{} {
	input.Model = new(models.ChangeUserEmail)
	input.Model.VerificationKey = models.GenerateSecretKey()
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *ChangeUserEmail) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil
}

// Validate is current model is valid
func (input *ChangeUserEmail) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Email == "" {
		result.AddFieldFailure("email", "Email is required.")
		return result
	}

	if len(input.Model.Email) > 200 {
		result.AddFieldFailure("email", "Email must have less than 200 characters.")
		return result
	}

	if user.Email == input.Model.Email {
		result.AddFieldFailure("email", "Choose a different email.")
		return result
	}

	messages := validate.Email(input.Model.Email)
	if len(messages) > 0 {
		result.AddFieldFailure("email", messages...)
		return result
	}

	existing, err := services.Users.GetByEmail(input.Model.Email)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	}
	if err == nil && existing.ID != user.ID {
		result.AddFieldFailure("email", "This email is already in use by someone else")
		return result
	}
	input.Model.Requestor = user
	return result
}
