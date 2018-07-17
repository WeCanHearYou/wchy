package validate

import (
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/img"
)

//ImageUpload validates given image upload
func ImageUpload(upload *models.ImageUpload, minWidth, minHeight int) *Result {
	messages := []string{}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		logo, err := img.Parse(upload.Upload.Content)
		if err != nil {
			if err == img.ErrNotSupported {
				messages = append(messages, "This file format not supported.")
			} else {
				return Error(err)
			}
		} else {
			if logo.Width < minWidth || logo.Height < minHeight {
				messages = append(messages, fmt.Sprintf("The image must have minimum dimensions of %dx%d pixels.", minWidth, minHeight))
			}

			if logo.Width != logo.Height {
				messages = append(messages, "The image must have an aspect ratio of 1:1.")
			}

			if logo.Size > 51200 {
				messages = append(messages, "The image size must be smaller than 100KB.")
			}
		}
	}

	if len(messages) > 0 {
		return Failed(messages)
	}
	return Success()
}