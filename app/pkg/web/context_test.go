package web_test

import (
	"crypto/tls"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

func newGetContext(rawurl string, headers map[string]string) *web.Context {
	u, _ := url.Parse(rawurl)
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", u.RequestURI(), nil)
	req.Host = u.Host

	if u.Scheme == "https" {
		req.TLS = &tls.ConnectionState{}
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	ctx := e.NewContext(res, req, nil)
	return &ctx
}

func newBodyContext(method string, params web.StringMap, body, contentType string) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, "/some/resource", strings.NewReader(body))
	req.Host = "demo.test.fider.io:3000"
	req.Header.Set("Content-Type", contentType)
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func TestContextID(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)

	Expect(ctx.ContextID()).IsNotEmpty()
	Expect(ctx.ContextID()).HasLen(32)
}

func TestBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)

	Expect(ctx.BaseURL()).Equals("http://demo.test.fider.io:3000")
}

func TestBaseURL_HTTPS(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("https://demo.test.fider.io:3000", nil)

	Expect(ctx.BaseURL()).Equals("https://demo.test.fider.io:3000")
}

func TestBaseURL_HTTPS_Proxy(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", map[string]string{
		"X-Forwarded-Proto": "https",
	})

	Expect(ctx.BaseURL()).Equals("https://demo.test.fider.io:3000")
}

func TestCurrentURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000/resource?id=23", nil)

	Expect(ctx.Request.URL.String()).Equals("http://demo.test.fider.io:3000/resource?id=23")
}

func TestTenantURL(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://login.test.fider.io:3000", nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://theavengers.test.fider.io:3000")
}

func TestTenantURL_WithCNAME(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "feedback.theavengers.com",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://feedback.theavengers.com:3000")
}

func TestTenantURL_SingleHostMode(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "single"

	ctx := newGetContext("http://demo.test.fider.io:3000", nil)
	tenant := &models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	}
	Expect(ctx.TenantBaseURL(tenant)).Equals("http://demo.test.fider.io:3000")
}

func TestGlobalAssetsURL_SingleHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "single"
	ctx := newGetContext("http://feedback.theavengers.com:3000", nil)
	Expect(ctx.GlobalAssetsURL("/assets/main.js")).Equals("http://feedback.theavengers.com:3000/assets/main.js")
	Expect(ctx.GlobalAssetsURL("/assets/main.css")).Equals("http://feedback.theavengers.com:3000/assets/main.css")

	env.Config.CDN.Host = "assets-fider.io"
	Expect(ctx.GlobalAssetsURL("/assets/main.js")).Equals("http://assets-fider.io/assets/main.js")
	Expect(ctx.GlobalAssetsURL("/assets/main.css")).Equals("http://assets-fider.io/assets/main.css")
}

func TestGlobalAssetsURL_MultiHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "multi"
	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)
	ctx.SetTenant(&models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "feedback.theavengers.com",
	})

	Expect(ctx.GlobalAssetsURL("/assets/main.js")).Equals("http://theavengers.test.fider.io:3000/assets/main.js")
	Expect(ctx.GlobalAssetsURL("/assets/main.css")).Equals("http://theavengers.test.fider.io:3000/assets/main.css")

	env.Config.CDN.Host = "assets-fider.io"
	Expect(ctx.GlobalAssetsURL("/assets/main.js")).Equals("http://cdn.assets-fider.io/assets/main.js")
	Expect(ctx.GlobalAssetsURL("/assets/main.css")).Equals("http://cdn.assets-fider.io/assets/main.css")
}

func TestTenantAssetsURL_SingleHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "single"
	ctx := newGetContext("http://feedback.theavengers.com:3000", nil)
	ctx.SetTenant(&models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
	})

	Expect(ctx.TenantAssetsURL("/assets/main.js")).Equals("http://feedback.theavengers.com:3000/assets/main.js")
	Expect(ctx.TenantAssetsURL("/assets/main.css")).Equals("http://feedback.theavengers.com:3000/assets/main.css")

	env.Config.CDN.Host = "assets-fider.io"
	Expect(ctx.TenantAssetsURL("/assets/main.js")).Equals("http://assets-fider.io/assets/main.js")
	Expect(ctx.TenantAssetsURL("/assets/main.css")).Equals("http://assets-fider.io/assets/main.css")
}

func TestTenantAssetsURL_MultiHostMode(t *testing.T) {
	RegisterT(t)

	env.Config.HostMode = "multi"
	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)
	ctx.SetTenant(&models.Tenant{
		ID:        1,
		Subdomain: "theavengers",
		CNAME:     "feedback.theavengers.com",
	})

	Expect(ctx.TenantAssetsURL("/assets/main.js")).Equals("http://theavengers.test.fider.io:3000/assets/main.js")
	Expect(ctx.TenantAssetsURL("/assets/main.css")).Equals("http://theavengers.test.fider.io:3000/assets/main.css")

	env.Config.CDN.Host = "assets-fider.io"
	Expect(ctx.TenantAssetsURL("/assets/main.js")).Equals("http://theavengers.assets-fider.io/assets/main.js")
	Expect(ctx.TenantAssetsURL("/assets/main.css")).Equals("http://theavengers.assets-fider.io/assets/main.css")
}

func TestCanonicalURL_SameDomain(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)

	ctx.SetCanonicalURL("")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000`)

	ctx.SetCanonicalURL("/some-url")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/some-url`)

	ctx.SetCanonicalURL("/some-other-url")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/some-other-url`)

	ctx.SetCanonicalURL("page-b/abc.html")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://theavengers.test.fider.io:3000/page-b/abc.html`)
}

func TestCanonicalURL_DifferentDomain(t *testing.T) {
	RegisterT(t)

	ctx := newGetContext("http://theavengers.test.fider.io:3000", nil)

	ctx.SetCanonicalURL("http://feedback.theavengers.com/some-url")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://feedback.theavengers.com/some-url`)

	ctx.SetCanonicalURL("")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://feedback.theavengers.com`)

	ctx.SetCanonicalURL("/some-other-url")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://feedback.theavengers.com/some-other-url`)

	ctx.SetCanonicalURL("page-b/abc.html")
	Expect(ctx.Get("Canonical-URL")).Equals(`http://feedback.theavengers.com/page-b/abc.html`)
}
