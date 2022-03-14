package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	mockUserResp := `{"message":"Hello World"}`

	ts := httptest.NewServer(SetupServer())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, actual %v", resp.StatusCode)
	}

	respData, _ := ioutil.ReadAll(resp.Body)
	if string(respData) != mockUserResp {
		t.Fatalf("Expected response body %v, actual %v", mockUserResp, string(respData))
	}
}
