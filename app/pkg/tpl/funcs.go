package tpl

import (
	"errors"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/crypto"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/microcosm-cc/bluemonday"
)

var strictHtmlPolicy = bluemonday.NewPolicy()

var templateFunctions = map[string]interface{}{
	"stripHtml": func(input string) string {
		return strictHtmlPolicy.Sanitize(input)
	},
	"html": func(input string) template.HTML {
		return template.HTML(input)
	},
	"md5": func(input string) string {
		return crypto.MD5(input)
	},
	"lower": func(input string) string {
		return strings.ToLower(input)
	},
	"upper": func(input string) string {
		return strings.ToUpper(input)
	},
	"translate": func(input string, params ...i18n.Params) string {
		return "This is overwritten later on..."
	},
	"markdown": func(input string) template.HTML {
		return markdown.Full(input)
	},
	"dict": func(values ...interface{}) map[string]interface{} {
		if len(values)%2 != 0 {
			panic(errors.New("invalid dictionary call"))
		}

		dict := make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			var key string
			switch v := values[i].(type) {
			case string:
				key = v
			default:
				panic(errors.New("invalid dictionary key"))
			}
			dict[key] = values[i+1]
		}
		return dict
	},
	"format": func(format string, date time.Time) string {
		return date.Format(format)
	},
	"quote": func(text string) string {
		return strconv.Quote(text)
	},
}
