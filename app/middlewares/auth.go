package middlewares

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

// IsAuthenticated blocks non-authenticated requests
func IsAuthenticated() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if !c.IsAuthenticated() {
				return c.Unauthorized()
			}
			return next(c)
		}
	}
}

// CheckAuthTokenOrigin blocks API-originated AuthToken requests a non-API resource
func CheckAuthTokenOrigin() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			claims := c.Claims()
			if claims != nil {
				if !c.Request.IsAPI() && claims.Origin == jwt.FiderClaimsOriginAPI {
					return c.Unauthorized()
				}
			}
			return next(c)
		}
	}
}

// IsAuthorized blocks non-authorized requests
func IsAuthorized(roles ...models.Role) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			user := c.User()
			for _, role := range roles {
				if user.Role == role {
					return next(c)
				}
			}
			return c.Unauthorized()
		}
	}
}
