package web_test

import (
	"net/url"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestReactRenderer_FileNotFound(t *testing.T) {
	RegisterT(t)

	r := web.NewReactRenderer("unknown.js")
	u, _ := url.Parse("https://github.com")
	html, err := r.Render(u, web.Map{})
	Expect(html).Equals("")
	Expect(err).IsNil()
}

func TestReactRenderer_RenderEmptyHomeHTML(t *testing.T) {
	RegisterT(t)

	r := web.NewReactRenderer("ssr.js")
	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"tenant":   &models.Tenant{},
		"settings": web.Map{},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).Equals(`<div id="c-header"><div class="c-env-info">Env:  | Compiler:  | Version:  | BuildTime: N/A |TenantID: 0 | </div><div class="c-menu"><div class="container"><a href="/" class="c-menu-item-title"><span></span></a><div class="c-menu-item-signin"><span>Sign in</span></div></div></div></div><div id="p-home" class="page container"><div class="row"><div class="l-welcome-col col-md-4"><div class="markdown-body welcome-message"><p>We&#39;d love to hear what you&#39;re thinking about. </p>
<p>What can we do better? This is the place for you to vote, discuss and share ideas.</p></div><form autoComplete="off" class="c-form"><div class="c-form-field"><div class="c-form-field-wrapper"><input type="text" id="input-title" tabindex="-1" maxLength="100" value="" placeholder="Enter your suggestion here..."/></div></div></form></div><div class="l-posts-col col-md-8"><div class="l-lonely center"><p><svg stroke="currentColor" fill="currentColor" stroke-width="0" viewBox="0 0 352 512" class="icon" height="1em" width="1em" xmlns="http://www.w3.org/2000/svg"><path d="M176 80c-52.94 0-96 43.06-96 96 0 8.84 7.16 16 16 16s16-7.16 16-16c0-35.3 28.72-64 64-64 8.84 0 16-7.16 16-16s-7.16-16-16-16zM96.06 459.17c0 3.15.93 6.22 2.68 8.84l24.51 36.84c2.97 4.46 7.97 7.14 13.32 7.14h78.85c5.36 0 10.36-2.68 13.32-7.14l24.51-36.84c1.74-2.62 2.67-5.7 2.68-8.84l.05-43.18H96.02l.04 43.18zM176 0C73.72 0 0 82.97 0 176c0 44.37 16.45 84.85 43.56 115.78 16.64 18.99 42.74 58.8 52.42 92.16v.06h48v-.12c-.01-4.77-.72-9.51-2.15-14.07-5.59-17.81-22.82-64.77-62.17-109.67-20.54-23.43-31.52-53.15-31.61-84.14-.2-73.64 59.67-128 127.95-128 70.58 0 128 57.42 128 128 0 30.97-11.24 60.85-31.65 84.14-39.11 44.61-56.42 91.47-62.1 109.46a47.507 47.507 0 0 0-2.22 14.3v.1h48v-.05c9.68-33.37 35.78-73.18 52.42-92.16C335.55 260.85 352 220.37 352 176 352 78.8 273.2 0 176 0z"></path></svg></p><p>It&#x27;s lonely out here. Start by sharing a suggestion!</p></div></div></div></div><div id="c-footer"><div class="container"><a class="l-powered" target="_blank" href="https://getfider.com/"><img src="https://getfider.com/images/logo-100x100.png" alt="Fider"/><span>Powered by Fider</span></a></div></div>`)
	Expect(err).IsNil()
}
