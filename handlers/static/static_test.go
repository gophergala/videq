package static

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCssFileLoad(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := "/resources/css/bootstrap.css"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	staticHandler := NewHandler("../../")

	staticHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatal(fmt.Sprintf("Server repled with status %v expected %v", resp.Code, http.StatusOK))
	}
}
