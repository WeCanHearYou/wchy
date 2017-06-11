package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

// JwtGetter gets JWT token from cookie and insert into context
func JwtGetter() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			var err error

			if cookie, err := c.Cookie("auth"); err == nil {
				if claims, err := jwt.Decode(cookie.Value); err == nil {
					services := c.Services()

					if user, err := services.Users.GetByID(claims.UserID); err == nil {
						if c.Tenant() != nil && user.Tenant.ID == c.Tenant().ID {
							c.SetUser(user)
							if !c.IsAjax() {
								if ids, err := services.Ideas.SupportedBy(user.ID); err == nil {
									c.AddRenderVar("supportedIdeas", ids)
								}
							}
						}
					}

				}
			}

			if err != nil {
				return err
			}

			return next(c)
		}
	}
}

// JwtSetter sets JWT token into cookie
func JwtSetter() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {

			query := c.Request().URL.Query()

			jwt := query.Get("jwt")
			if jwt != "" {
				c.SetCookie(&http.Cookie{
					Name:     "auth",
					Value:    jwt,
					HttpOnly: true,
					Path:     "/",
					Expires:  time.Now().Add(365 * 24 * time.Hour),
				})

				query.Del("jwt")

				url := c.BaseURL() + c.Request().URL.Path
				querystring := query.Encode()
				if querystring != "" {
					url += "?" + querystring
				}

				return c.Redirect(http.StatusTemporaryRedirect, url)
			}

			return next(c)
		}
	}
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
