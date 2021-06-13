package web_test

import (
	"net/url"
	"testing"

	"github.com/getfider/fider/app/models/entity"
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
		"tenant":   &entity.Tenant{},
		"settings": web.Map{
			"locale": "en",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<input type="text" class="c-input" id="input-title"`)
	Expect(html).ContainsSubstring(`What can we do better? This is the place for you to vote, discuss and share ideas.`)
	Expect(html).ContainsSubstring(`No posts have been created yet.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}

func TestReactRenderer_RenderEmptyHomeHTML_Portuguese(t *testing.T) {
	RegisterT(t)

	r := web.NewReactRenderer("ssr.js")
	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"tenant":   &entity.Tenant{},
		"settings": web.Map{
			"locale": "pt-BR",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<input type="text" class="c-input" id="input-title"`)
	Expect(html).ContainsSubstring(`What can we do better? This is the place for you to vote, discuss and share ideas.`)
	Expect(html).ContainsSubstring(`Nenhum post foi criado ainda.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}
