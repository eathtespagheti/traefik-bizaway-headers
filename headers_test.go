package headerrules

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const PLUGIN_NAME = "heders-plugin"
const DEFAULT_URL = "http://localhost"
const TEST_HEADER = "X-Bizaway-Test"

func assertReqHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}

func assertResHeader(t *testing.T, res *http.Response, key, expected string) {
	t.Helper()

	if res.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", res.Header.Get(key))
	}
}

var mock_server_response = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set(TEST_HEADER, "Undesired-Value")
	rw.WriteHeader(http.StatusOK)
})


func TestRequestHeaders(t *testing.T) {
	cfg := CreateConfig()
	cfg.Headers.Request["X-Host"] = "[[.Host]]"
	cfg.Headers.Request["X-Method"] = "[[.Method]]"
	cfg.Headers.Request["X-URL"] = "[[.URL]]"
	cfg.Headers.Request["X-Demo"] = "test"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { /* Dummy function */ })

	handler, err := New(ctx, next, cfg, PLUGIN_NAME)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DEFAULT_URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertReqHeader(t, req, "X-Host", "localhost")
	assertReqHeader(t, req, "X-URL", DEFAULT_URL)
	assertReqHeader(t, req, "X-Method", "GET")
	assertReqHeader(t, req, "X-Demo", "test")
}

func TestResponseHeaders(t *testing.T) {
	cfg := CreateConfig()
	cfg.Headers.Response[TEST_HEADER] = TEST_HEADER
	ctx := context.Background()

	handler, err := New(ctx, mock_server_response, cfg, PLUGIN_NAME)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DEFAULT_URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	assertResHeader(t, recorder.Result(), TEST_HEADER, TEST_HEADER);
}

func TestNoHeadersConfigured(t *testing.T) {
	cfg := CreateConfig()
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { /* Dummy function */ })

	_, err := New(ctx, next, cfg, PLUGIN_NAME)
	if err == nil {
		t.Fatal("expected error when no headers are configured")
	}
	if err.Error() != "at least one request or response header must be configured" {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestRequestResponseChain(t *testing.T) {
	/*
	- Setup a headers replacement for both Request and Response
	- Create the request
	- Process it with mock_server_response
	- Ensure that the headers have been correctly replaced on both request and response
	*/
	cfg := CreateConfig()
	cfg.Headers.Request[TEST_HEADER] = TEST_HEADER
	cfg.Headers.Response[TEST_HEADER] = TEST_HEADER
	// Setup host and allow origin update
	cfg.Headers.Request["X-Original-Host"] = "[[.Host]]"
	cfg.Headers.Request["Host"] = "custom.endpoint"
	cfg.Headers.Response["Access-Control-Allow-Origin"] = "https://[[.Request.Header.Get \"X-Original-Host\"]]"


	ctx := context.Background()

	handler, err := New(ctx, mock_server_response, cfg, PLUGIN_NAME)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://app.domain.name", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	assertReqHeader(t, req, TEST_HEADER, TEST_HEADER)
	assertResHeader(t, recorder.Result(), TEST_HEADER, TEST_HEADER)
	assertReqHeader(t, req, "X-Original-Host", "app.domain.name")
	assertReqHeader(t, req, "Host", "custom.endpoint")
	assertResHeader(t, recorder.Result(), "Access-Control-Allow-Origin", "https://app.domain.name")
}




