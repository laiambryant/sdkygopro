package endpoint

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/laiambryant/sdkyGOprodeck/client"
)

type fakeHTTP struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (f *fakeHTTP) Do(req *http.Request) (*http.Response, error) { return f.fn(req) }

func TestNewEndpoint(t *testing.T) {
	c := client.NewHTTPClient(&fakeHTTP{fn: func(req *http.Request) (*http.Response, error) { return client.NewMockResponse(200, `{}`), nil }}, client.WithBaseURL("http://example"))
	e := New[struct{}](c)
	if e.Client != c {
		t.Fatalf("expected client to be set")
	}
}

func TestFetchSuccessAndDecodeErrorAndRequestError(t *testing.T) {
	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	c := client.NewHTTPClient(&fakeHTTP{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "http://example/cardinfo.php?name=test" {
			t.Fatalf("unexpected url: %s", req.URL.String())
		}
		return client.NewMockResponse(200, `{"id":123,"name":"bob"}`), nil
	}}, client.WithBaseURL("http://example"))
	e := New[Item](c)
	it, err := e.Fetch(context.Background(), "/cardinfo.php?name=test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if it.ID != 123 || it.Name != "bob" {
		t.Fatalf("unexpected item: %#v", it)
	}

	c2 := client.NewHTTPClient(&fakeHTTP{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(200, `not-json`), nil
	}}, client.WithBaseURL("http://example"))
	e2 := New[Item](c2)
	_, err = e2.Fetch(context.Background(), "/cardinfo.php")
	if err == nil {
		t.Fatalf("expected decode error")
	}
	var derr *DecodeError
	if !errors.As(err, &derr) {
		t.Fatalf("expected DecodeError, got %T", err)
	}
	if !strings.Contains(derr.Error(), "/cardinfo.php") {
		t.Fatalf("error message should contain resource: %v", derr.Error())
	}
	if derr.Unwrap() == nil {
		t.Fatalf("expected wrapped error")
	}

	c3 := client.NewHTTPClient(&fakeHTTP{fn: func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	}}, client.WithBaseURL("http://example"))
	e3 := New[Item](c3)
	_, err = e3.Fetch(context.Background(), "/cardinfo.php")
	if err == nil {
		t.Fatalf("expected request error")
	}
	var re *client.RequestError
	if !errors.As(err, &re) {
		t.Fatalf("expected client.RequestError, got %T", err)
	}
}
