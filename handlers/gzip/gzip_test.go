package gzip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

func TestGzipPassHandler(t *testing.T) {

	gzipHandlerWraper := NewHandler(&testHandler{})

	ts := httptest.NewServer(gzipHandlerWraper)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header = http.Header{"Accept-Encoding": {"gzip"}}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		t.Fatal(fmt.Sprintf("Gzip header expected in response header"))
	}
}

func TestContentNotGziped(t *testing.T) {

	gzipHandlerWraper := NewHandler(&testHandler{})

	ts := httptest.NewServer(gzipHandlerWraper)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// client dose not accepts gzip on purpuse
	//req.Header = http.Header{"Accept-Encoding": {"gzip"}}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		t.Fatal(fmt.Sprintf("Gzip header NOT expected in response header"))
	}
}
