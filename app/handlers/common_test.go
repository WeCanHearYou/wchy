package handlers_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestHealthHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.Health())

	Expect(code).Equals(http.StatusOK)
}

func TestPageHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.Page("The Title", "The Description"))

	Expect(code).Equals(http.StatusOK)
}

func TestLegalPageHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.LegalPage("Terms of Service", "terms.md"))

	Expect(code).Equals(http.StatusOK)
}

func TestLegalPageHandler_Invalid(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.LegalPage("Some Page", "somepage.md"))

	Expect(code).Equals(http.StatusNotFound)
}

func TestRobotsTXT(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.Execute(handlers.RobotsTXT())
	Expect(code).Equals(http.StatusOK)
}

func TestSitemap(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io:3000/sitemap.xml").
		Execute(handlers.Sitemap())

	bytes, _ := ioutil.ReadAll(response.Body)
	Expect(code).Equals(http.StatusOK)
	Expect(string(bytes)).Equals(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url> <loc>http://demo.test.fider.io:3000</loc> </url></urlset>`)
}

func TestSitemap_WithPosts(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Posts.SetCurrentUser(mock.AryaStark)
	services.Posts.Add("My new idea 1", "")
	services.Posts.SetCurrentUser(mock.AryaStark)
	services.Posts.Add("The other idea", "")

	code, response := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io:3000/sitemap.xml").
		Execute(handlers.Sitemap())

	bytes, _ := ioutil.ReadAll(response.Body)
	Expect(code).Equals(http.StatusOK)
	Expect(string(bytes)).Equals(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url> <loc>http://demo.test.fider.io:3000</loc> </url><url> <loc>http://demo.test.fider.io:3000/posts/1/my-new-idea-1</loc> </url><url> <loc>http://demo.test.fider.io:3000/posts/2/the-other-idea</loc> </url></urlset>`)
}

func TestSitemap_PrivateTenant_WithPosts(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Posts.SetCurrentUser(mock.AryaStark)
	services.Posts.Add("My new idea 1", "")
	services.Posts.SetCurrentUser(mock.AryaStark)
	services.Posts.Add("The other idea", "")

	mock.DemoTenant.IsPrivate = true

	code, _ := server.
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io:3000/sitemap.xml").
		Execute(handlers.Sitemap())

	Expect(code).Equals(http.StatusNotFound)
}
