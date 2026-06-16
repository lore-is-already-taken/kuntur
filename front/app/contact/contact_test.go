package contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// stubRoundTripper is a configurable http.RoundTripper for tests. It records
// the most recent request body and returns the configured response.
type stubRoundTripper struct {
	statusCode int
	body       string
	err        error
	gotBody    []byte
	gotURL     string
	gotMethod  string
}

func (s *stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if s.err != nil {
		return nil, s.err
	}
	if r.Body != nil {
		b, _ := io.ReadAll(r.Body)
		s.gotBody = b
	}
	s.gotURL = r.URL.String()
	s.gotMethod = r.Method
	return &http.Response{
		StatusCode: s.statusCode,
		Body:       io.NopCloser(strings.NewReader(s.body)),
		Header:     make(http.Header),
	}, nil
}

// templateWithSuccess returns a minimal template that exposes the layout
// block with a body that includes the Success flag, so tests can assert
// whether the GET handler toggled it correctly.
func templateWithSuccess(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("layout").Parse(`{{ define "layout" }}success={{ .Success }}{{ end }}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return tmpl
}

func TestContactGet_DefaultSuccessFalse(t *testing.T) {
	tmpl := templateWithSuccess(t)
	h := Get(tmpl)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/contacto", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "success=false") {
		t.Fatalf("body: got %q, want substring success=false", rec.Body.String())
	}
}

func TestContactGet_OK1TogglesSuccessTrue(t *testing.T) {
	tmpl := templateWithSuccess(t)
	h := Get(tmpl)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/contacto?ok=1", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "success=true") {
		t.Fatalf("body: got %q, want substring success=true", rec.Body.String())
	}
}

func TestContactPost_HappyPathRedirectsToOK1(t *testing.T) {
	stub := &stubRoundTripper{statusCode: 200, body: "{}"}
	client := &http.Client{Transport: stub, Timeout: 5 * time.Second}
	h := New(client, "http://upstream.test/contacto")
	rec := httptest.NewRecorder()
	form := "name=Ada&email=a%40b.c&message=hi"
	req := httptest.NewRequest(http.MethodPost, "/contacto", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: got %d, want 303", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/contacto?ok=1" {
		t.Fatalf("location: got %q, want /contacto?ok=1", loc)
	}
	if stub.gotMethod != http.MethodPost {
		t.Fatalf("upstream method: got %q, want POST", stub.gotMethod)
	}
	if !strings.Contains(stub.gotURL, "http://upstream.test/contacto") {
		t.Fatalf("upstream url: got %q, want contains upstream URL", stub.gotURL)
	}
	var p Payload
	if err := json.Unmarshal(stub.gotBody, &p); err != nil {
		t.Fatalf("upstream body: not JSON: %v (body=%q)", err, stub.gotBody)
	}
	if p.Name != "Ada" || p.Email != "a@b.c" || p.Message != "hi" {
		t.Fatalf("upstream payload: got %+v", p)
	}
}

func TestContactPost_Upstream4xxReturns502(t *testing.T) {
	stub := &stubRoundTripper{statusCode: 400, body: "bad request"}
	client := &http.Client{Transport: stub, Timeout: 5 * time.Second}
	h := New(client, "http://upstream.test/contacto")
	rec := httptest.NewRecorder()
	form := "name=Ada&email=a%40b.c&message=hi"
	req := httptest.NewRequest(http.MethodPost, "/contacto", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d, want 502", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "upstream error") {
		t.Fatalf("body: got %q, want substring 'upstream error'", rec.Body.String())
	}
}

func TestContactPost_Upstream5xxReturns502(t *testing.T) {
	stub := &stubRoundTripper{statusCode: 502, body: "bad gateway"}
	client := &http.Client{Transport: stub, Timeout: 5 * time.Second}
	h := New(client, "http://upstream.test/contacto")
	rec := httptest.NewRecorder()
	form := "name=Ada&email=a%40b.c&message=hi"
	req := httptest.NewRequest(http.MethodPost, "/contacto", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d, want 502", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "upstream error") {
		t.Fatalf("body: got %q, want substring 'upstream error'", rec.Body.String())
	}
}

func TestContactPost_NetworkErrorReturns502(t *testing.T) {
	stub := &stubRoundTripper{err: errors.New("connection refused")}
	client := &http.Client{Transport: stub, Timeout: 5 * time.Second}
	h := New(client, "http://upstream.test/contacto")
	rec := httptest.NewRecorder()
	form := "name=Ada&email=a%40b.c&message=hi"
	req := httptest.NewRequest(http.MethodPost, "/contacto", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: got %d, want 502", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "service unavailable") {
		t.Fatalf("body: got %q, want substring 'service unavailable'", rec.Body.String())
	}
}

func TestContactPost_MalformedFormReturns400(t *testing.T) {
	client := &http.Client{Transport: &stubRoundTripper{statusCode: 200}, Timeout: 5 * time.Second}
	h := New(client, "http://upstream.test/contacto")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/contacto", bytes.NewReader([]byte("not a form")))
	req.Header.Set("Content-Type", "this is not a valid content type \x00with binary")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid form") {
		t.Fatalf("body: got %q, want substring 'invalid form'", rec.Body.String())
	}
}
