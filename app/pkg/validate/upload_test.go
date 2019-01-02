package validate_test

import (
	"io/ioutil"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestValidateImageUpload(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		fileName string
		count    int
	}{
		{"/app/pkg/img/testdata/logo1.png", 0},
		{"/app/pkg/img/testdata/logo2.jpg", 2},
		{"/app/pkg/img/testdata/logo3.gif", 1},
		{"/app/pkg/img/testdata/logo4.png", 1},
		{"/app/pkg/img/testdata/logo5.png", 0},
		{"/README.md", 1},
		{"/app/pkg/img/testdata/favicon.ico", 1},
	}

	for _, testCase := range testCases {
		img, _ := ioutil.ReadFile(env.Path(testCase.fileName))

		upload := &models.ImageUpload{
			Upload: &models.ImageUploadData{
				Content: img,
			},
		}
		messages, err := validate.ImageUpload(upload, validate.ImageUploadOpts{
			MinHeight:    200,
			MinWidth:     200,
			MaxKilobytes: 100,
			ExactRatio:   true,
		})
		Expect(messages).HasLen(testCase.count)
		Expect(err).IsNil()
	}
}

func TestValidateImageUpload_ExactRatio(t *testing.T) {
	RegisterT(t)

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/img/testdata/logo3-200w.gif"))
	opts := validate.ImageUploadOpts{
		IsRequired:   false,
		MaxKilobytes: 200,
	}

	upload := &models.ImageUpload{
		Upload: &models.ImageUploadData{
			Content: img,
		},
	}
	opts.ExactRatio = true
	messages, err := validate.ImageUpload(upload, opts)
	Expect(messages).HasLen(1)
	Expect(err).IsNil()

	opts.ExactRatio = false
	messages, err = validate.ImageUpload(upload, opts)
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}

func TestValidateImageUpload_Nil(t *testing.T) {
	RegisterT(t)

	messages, err := validate.ImageUpload(nil, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    200,
		MinWidth:     200,
		MaxKilobytes: 50,
		ExactRatio:   true,
	})
	Expect(messages).HasLen(0)
	Expect(err).IsNil()

	messages, err = validate.ImageUpload(&models.ImageUpload{}, validate.ImageUploadOpts{
		IsRequired:   false,
		MinHeight:    200,
		MinWidth:     200,
		MaxKilobytes: 50,
		ExactRatio:   true,
	})
	Expect(messages).HasLen(0)
	Expect(err).IsNil()
}

func TestValidateImageUpload_Required(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		upload *models.ImageUpload
		count  int
	}{
		{nil, 1},
		{&models.ImageUpload{}, 1},
		{&models.ImageUpload{
			BlobKey: "some-file.png",
			Remove:  true,
		}, 1},
	}

	for _, testCase := range testCases {
		messages, err := validate.ImageUpload(testCase.upload, validate.ImageUploadOpts{
			IsRequired:   true,
			MinHeight:    200,
			MinWidth:     200,
			MaxKilobytes: 50,
			ExactRatio:   true,
		})
		Expect(messages).HasLen(testCase.count)
		Expect(err).IsNil()
	}
}
