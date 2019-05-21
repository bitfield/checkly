package checkly

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
)

func assertFormParamPresent(t *testing.T, form url.Values, param string) {
	value := form.Get(param)
	if value == "" {
		t.Errorf("want %q parameter, got none", param)
	}
}
func TestCreate(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("want POST request, got %q", r.Method)
		}
		wantURL := "/v1/checks"
		if r.URL.EscapedPath() != wantURL {
			t.Errorf("want %q, got %q", wantURL, r.URL.EscapedPath())
		}
		r.ParseForm()
		wantParams := []string{"name", "checkType", "activated"}
		for _, p := range wantParams {
			assertFormParamPresent(t, r.Form, p)
		}
		w.WriteHeader(http.StatusOK)
		data, err := os.Open("testdata/Create.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client := NewClient("dummy")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	wantID := "73d29e72-6540-4bb5-967e-e07fa2c9465e"
	check := Check{
		Name:      "test",
		Type:      TypeBrowser,
		Activated: true,
	}
	gotID, err := client.Create(check)
	if err != nil {
		t.Fatal(err)
	}
	if gotID != wantID {
		t.Errorf("want %q, got %q", wantID, gotID)
	}
}

const idFormat = `[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}`

var idRE = regexp.MustCompile(idFormat)

func TestDelete(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("want DELETE request, got %q", r.Method)
		}
		wantURLPrefix := "/v1/checks"
		if !strings.HasPrefix(r.URL.EscapedPath(), wantURLPrefix) {
			t.Errorf("want URL prefix %q, got %q", wantURLPrefix, r.URL.EscapedPath())
		}
		ID := path.Base(r.URL.String())
		if !idRE.MatchString(ID) {
			t.Errorf("malformed ID %q (should match %q)", ID, idFormat)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	client := NewClient("dummy")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	err := client.Delete("73d29e72-6540-4bb5-967e-e07fa2c9465e")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("want GET request, got %q", r.Method)
		}
		wantURLPrefix := "/v1/checks"
		if !strings.HasPrefix(r.URL.EscapedPath(), wantURLPrefix) {
			t.Errorf("want URL prefix %q, got %q", wantURLPrefix, r.URL.EscapedPath())
		}
		ID := path.Base(r.URL.String())
		if !idRE.MatchString(ID) {
			t.Errorf("malformed ID %q (should match %q)", ID, idFormat)
		}
		w.WriteHeader(http.StatusOK)
		data, err := os.Open("testdata/Create.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	client := NewClient("dummy")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	check, err := client.Get("73d29e72-6540-4bb5-967e-e07fa2c9465e")
	if err != nil {
		t.Fatal(err)
	}
	wantURL := "http://example.com"
	if check.Request.URL != wantURL {
		t.Errorf("want URL %q, got %q", wantURL, check.Request.URL)
	}
}
