package mux

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func serve(code int) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(code)
	}
}

func TestMux(t *testing.T) {
	testing.Short()
	mux := New()

	routes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"/test", serve(200)},
		{"/testa", serve(201)},
		{"/test/long", serve(202)},
		{"/test/long/", serve(203)},
	}

	for _, r := range routes {
		mux.Handle(r.pattern, r.handler)
	}

	tc := map[string]struct {
		path    string
		pattern string
		code    int
	}{
		"longest path":                       {"/test/long/is/good", "/test/long/", 203},
		"longest path with trailing slash":   {"/test/long/", "/test/long/", 203},
		"longest path without traling slash": {"/test/long", "/test/long", 202},
		"shortest path":                      {"/test", "/test", 200},
		"variant a shortest path":            {"/testa", "/testa", 201},
		"no partial match":                   {"/tes", "", 404},
	}

	for name, test := range tc {
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: test.path,
			},
		}
		h, pattern := mux.Handler(req)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if pattern != test.pattern || rr.Code != test.code {
			t.Errorf("%s : %s = %d, %q, want %d, %q", name, test.path, rr.Code, pattern, test.code, test.pattern)
		}
	}
}

func TestMuxPanicDuplicatePath(t *testing.T) {
	expected := "mux: duplicate path for /test"
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
		} else if r != expected {
			t.Errorf("Expected panic with message `%s` but got `%s`", expected, r)
		}
	}()

	mux := New()
	mux.HandleFunc("/test", serve(200))
	mux.HandleFunc("/test", serve(200))
}

func TestMuxPanicNilHandler(t *testing.T) {
	expected := "mux: nil handler"
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
		} else if r != expected {
			t.Errorf("Expected panic with message `%s` but got `%s`", expected, r)
		}
	}()

	mux := New()
	mux.HandleFunc("/test", nil)
}
