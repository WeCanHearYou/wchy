package handlers

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/img"
	"github.com/getfider/fider/app/pkg/log"
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
					url := fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=%d&d=404", crypto.MD5(strings.ToLower(user.Email)), size)
					c.Logger().Debugf("Requesting gravatar: @{GravatarURL}", log.Props{
						"GravatarURL": url,
					})
					resp, err := http.Get(url)
					if err == nil {
						defer resp.Body.Close()

						if resp.StatusCode == http.StatusOK {
							bytes, err := ioutil.ReadAll(resp.Body)
							if err == nil {
								return c.Image(http.DetectContentType(bytes), bytes)
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

		return c.Image("image/png", buf.Bytes())
	}
}

//Favicon returns the Fider favicon by given size
func Favicon() web.HandlerFunc {
	return func(c web.Context) error {
		var (
			bytes       []byte
			err         error
			contentType string
		)

		bkey := c.Param("bkey")
		if bkey != "" {
			logo, err := c.Services().Blobs.Get(bkey)
			if err != nil {
				return c.Failure(err)
			}
			bytes = logo.Object
			contentType = logo.ContentType
		} else {
			bytes, err = ioutil.ReadFile(env.Path("favicon.png"))
			contentType = "image/png"
			if err != nil {
				return c.Failure(err)
			}
		}

		size, err := c.ParamAsInt("size")
		if err != nil {
			return c.NotFound()
		}

		opts := []img.ImageOperation{
			img.Padding(10),
			img.Resize(size),
		}

		if c.QueryParam("bg") != "" {
			opts = append(opts, img.ChangeBackground(color.White))
		}

		bytes, err = img.Apply(bytes, opts...)
		if err != nil {
			return c.Failure(err)
		}

		return c.Image(contentType, bytes)
	}
}

//ViewUploadedImage returns any uploaded image by given ID and size
func ViewUploadedImage() web.HandlerFunc {
	return func(c web.Context) error {
		bkey := c.Param("bkey")

		size, err := c.ParamAsInt("size")
		if err != nil {
			return c.NotFound()
		}

		logo, err := c.Services().Blobs.Get(bkey)
		if err != nil {
			return c.Failure(err)
		}

		bytes, err := img.Apply(logo.Object, img.Resize(size))
		if err != nil {
			return c.Failure(err)
		}

		return c.Image(logo.ContentType, bytes)
	}
}
