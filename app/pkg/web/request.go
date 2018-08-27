package web

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/errors"
)

//Request wraps the http request object
type Request struct {
	instance      *http.Request
	Method        string
	ClientIP      string
	ContentLength int64
	Body          string
	IsSecure      bool
	URL           *url.URL
}

// WrapRequest returns Fider wrapper of HTTP Request
func WrapRequest(request *http.Request) Request {
	protocol := "http"
	if request.TLS != nil || request.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}

	host := request.Host
	if request.Header.Get("X-Forwarded-Host") != "" {
		host = request.Header.Get("X-Forwarded-Host")
	}

	fullURL := protocol + "://" + host + request.RequestURI
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse url '%s'", fullURL))
	}

	var bodyBytes []byte
	if request.ContentLength > 0 {
		bodyBytes, err = ioutil.ReadAll(request.Body)
		if err != nil {
			panic(errors.Wrap(err, "failed to read body").Error())
		}
	}

	clientIP := request.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = strings.Split(request.RemoteAddr, ":")[0]
		if clientIP == "" {
			clientIP = "N/A"
		}
	} else {
		clientIP = strings.Split(clientIP, ",")[0]
	}

	return Request{
		instance:      request,
		Method:        request.Method,
		ClientIP:      strings.TrimSpace(clientIP),
		ContentLength: request.ContentLength,
		Body:          string(bodyBytes),
		URL:           u,
		IsSecure:      protocol == "https",
	}
}

// GetHeader returns the value of HTTP header from given key
func (r *Request) GetHeader(key string) string {
	return r.instance.Header.Get(key)
}

// SetHeader updates the value of HTTP header of given key
func (r *Request) SetHeader(key, value string) {
	r.instance.Header.Set(key, value)
}

// Cookie returns the named cookie provided in the request.
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	cookie, err := r.instance.Cookie(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get '%s' cookie", name)
	}
	return cookie, nil
}

// AddCookie adds a cookie
func (r *Request) AddCookie(cookie *http.Cookie) {
	r.instance.AddCookie(cookie)
}

// IsAPI returns true if its a request for an API resource
func (r *Request) IsAPI() bool {
	return strings.HasPrefix(r.URL.Path, "/api/")
}
