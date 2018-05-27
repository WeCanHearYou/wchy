package handlers

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/img"
	"github.com/getfider/fider/app/pkg/md5"

	"github.com/getfider/fider/app/pkg/web"
	"github.com/goenning/letteravatar"
)

//Avatar returns a gravatar picture of fallsback to letter avatar based on name
func Avatar() web.HandlerFunc {
	return func(c web.Context) error {
		name := c.Param("name")
		size, _ := c.ParamAsInt("size")
		id, err := c.ParamAsInt("id")

		if err == nil && id > 0 {
			user, err := c.Services().Users.GetByID(id)
			if err == nil && user.Tenant.ID == c.Tenant().ID {
				if user.Email != "" {
					url := fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=%d&d=404", md5.Hash(user.Email), size)
					c.Logger().Debugf("Requesting gravatar: %s", url)
					resp, err := http.Get(url)
					if err == nil {
						defer resp.Body.Close()

						if resp.StatusCode == http.StatusOK {
							bytes, err := ioutil.ReadAll(resp.Body)
							if err == nil {
								return c.Blob(http.StatusOK, http.DetectContentType(bytes), bytes)
							}
						}
					}
				}
			}
		}

		img, err := letteravatar.Draw(size, strings.ToUpper(letteravatar.Extract(name)), &letteravatar.Options{
			PaletteKey: fmt.Sprintf("%d:%s", id, name),
		})
		if err != nil {
			return c.Failure(err)
		}

		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			return c.Failure(err)
		}

		return c.Blob(http.StatusOK, "image/png", buf.Bytes())
	}
}

//Logo returns tenant logo by its ID on a given size
func Logo() web.HandlerFunc {
	return func(c web.Context) error {
		id, err := c.ParamAsInt("id")
		if err != nil {
			return c.Failure(err)
		}

		size, err := c.ParamAsInt("size")
		if err != nil {
			return c.Failure(err)
		}

		logo, err := c.Services().Tenants.GetLogo(id)
		if err != nil {
			return c.Failure(err)
		}

		bytes, err := img.Resize(logo.Content, size)
		if err != nil {
			return c.Failure(err)
		}

		return c.Blob(http.StatusOK, logo.ContentType, bytes)
	}
}
