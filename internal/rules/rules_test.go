package rules

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHeaderRuleHeaderChangeConditions(t *testing.T) {
	rule := HeaderRule{}
	if !rule.headerChangeConditions() {
		t.Error("HeaderChangeConditions should return true by default")
	}
}

func TestHeaderRuleInitSourceHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:  req,
		Response: res,
		Source:   "request",
	}
	rule.initSourceHeader()
	if rule.sourceHeaders == nil || rule.sourceHeaders != &req.Header {
		t.Error("Expected sourceHeader to be initialized to request header")
	}

	rule.Source = "response"
	rule.initSourceHeader()
	if rule.sourceHeaders == nil || rule.sourceHeaders != &res.Header {
		t.Error("Expected sourceHeader to be initialized to response header")
	}

	rule.Source = "invalid"
	rule.initSourceHeader()
	if rule.sourceHeaders == nil {
		t.Error("Expected sourceHeader to be nil")
	}
}

func TestHeaderRuleInitDestinationHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:     req,
		Response:    res,
		Destination: "request",
	}
	rule.initDestinationHeader()
	if rule.destinationHeaders == nil || rule.destinationHeaders != &req.Header {
		t.Error("Expected destinationHeader to be initialized to request header")
	}

	rule.Destination = "response"
	rule.initDestinationHeader()
	if rule.destinationHeaders == nil || rule.destinationHeaders != &res.Header {
		t.Error("Expected destinationHeader to be initialized to response header")
	}
	rule.Destination = "invalid"
	rule.initDestinationHeader()
	if rule.destinationHeaders == nil {
		t.Error("Expected destinationHeader to be nil")
	}
}

func TestHeaderRuleValidate(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:            req,
		Response:           res,
		destinationHeaders: &req.Header,
		sourceHeaders:      &res.Header,
	}
	valid, err := rule.Validate()
	if !valid || err != nil {
		t.Errorf("Expected Validate to return true and nil error, got: %v, %v", valid, err)
	}

	rule.sourceHeaders = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid source header" {
		t.Errorf("Expected error message 'invalid source header', got: %v", err)
	}

	rule.sourceHeaders = &res.Header
	rule.destinationHeaders = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid destination header" {
		t.Errorf("Expected error message 'invalid destination header', got: %v", err)
	}
}

func TestStringHeaderRuleGetHeaderValue(t *testing.T) {
	rule := StringHeaderRule{
		HeaderRule: HeaderRule{
			Value: "testHeader",
		},
	}
	if rule.GetHeaderValue() != "testHeader" {
		t.Errorf("Expected GetHeaderValue to return 'testHeader', got: %s", rule.GetHeaderValue())
	}
}

func TestStringHeaderRuleValidate(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	rule := StringHeaderRule{
		HeaderRule: HeaderRule{
			destinationHeaders: &req.Header,
		},
	}

	valid, err := rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid source header" {
		t.Errorf("Expected error message 'invalid source header', got: %v", err)
	}

	rule.sourceHeaders = rule.destinationHeaders
	valid, err = rule.Validate()
	if !valid || err != nil {
		t.Errorf("Expected Validate to return true and nil error, got: %v, %v", valid, err)
	}

	rule.destinationHeaders = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid destination header" {
		t.Errorf("Expected error message 'invalid destination header', got: %v", err)
	}
}

func TestNewStringHeaderRule(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule, _ := NewStringHeaderRule("Test-Header", "testValue", "request", req, res, nil, false)

	if rule.Header != "Test-Header" {
		t.Errorf("Expected Header to be 'Test-Header', got: %s", rule.Header)
	}
	if rule.Value != "testValue" {
		t.Errorf("Expected Value to be 'testValue', got: %s", rule.Value)
	}
	if rule.Destination != "request" {
		t.Errorf("Expected Destination to be 'request', got: %s", rule.Destination)
	}
	if rule.Request != req {
		t.Error("Expected Request to be the same")
	}
	if rule.Response != res {
		t.Error("Expected Response to be the same")
	}
	if rule.destinationHeaders != &req.Header {
		t.Error("Expected destinationHeader to be initialized to request header")
	}
}

func TestNewStringHeaderRuleInvalid(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	_, err := NewStringHeaderRule("Test-Header", "testValue", "invalid", req, res, nil, false)
	// Expect to get a error
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStringHeaderRuleSetValue(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule, _ := NewStringHeaderRule("Test-Header", "testValue", "request", req, res, nil, false)

	rule.SetHeader()
	// Check that the value on the header has been correctly set
	if req.Header.Get(rule.Header) != rule.Value {
		t.Errorf("Expected header value to be '%s', got: %s", rule.Value, req.Header.Get(rule.Header))
	}
}

func TestStringHeaderRuleTemplate(t *testing.T) {
	req := &http.Request{Header: http.Header{}, URL: &url.URL{Host: "test.com"}}
	res := &http.Response{Header: http.Header{}}
	rule, _ := NewStringHeaderRule("Test-Header", "[[.Request.URL.Host]]", "request", req, res, nil, true)

	rule.SetHeader()
	// Check that the value on the header has been correctly set
	if req.Header.Get(rule.Header) != req.URL.Host {
		t.Errorf("Expected header value to be '%s', got: %s", req.URL.Host, req.Header.Get(rule.Header))
	}
}


func TestCopyHeaderRuleGetHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{"Test-Header": []string{"testValue"}}}
	rule := CopyHeaderRule{
		HeaderRule: HeaderRule{
			Header:        "Test-Header",
			sourceHeaders: &req.Header,
		},
	}
	if rule.GetHeaderValue() != "testValue" {
		t.Errorf("Expected GetHeaderValue to return 'testValue', got: %s", rule.GetHeaderValue())
	}
}

func TestNewCopyHeaderRule(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule, _ := NewCopyHeaderRule("Test-Header", "testValue", "request", "response", req, res)

	if rule.Header != "Test-Header" {
		t.Errorf("Expected Header to be 'Test-Header', got: %s", rule.Header)
	}
	if rule.Value != "testValue" {
		t.Errorf("Expected Value to be 'testValue', got: %s", rule.Value)
	}
	if rule.Source != "request" {
		t.Errorf("Expected Source to be 'request', got: %s", rule.Source)
	}
	if rule.Destination != "response" {
		t.Errorf("Expected Destination to be 'response', got: %s", rule.Destination)
	}
	if rule.Request != req {
		t.Error("Expected Request to be the same")
	}
	if rule.Response != res {
		t.Error("Expected Response to be the same")
	}
	if rule.sourceHeaders != &req.Header {
		t.Error("Expected sourceHeader to be initialized to request header")
	}
	if rule.destinationHeaders != &res.Header {
		t.Error("Expected destinationHeader to be initialized to response header")
	}
}

func TestNewCopyHeaderRuleInvalid(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	_, err := NewCopyHeaderRule("Test-Header", "testValue", "invalid", "response", req, res)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestHeaderRuleGetHeaderValue(t *testing.T) {
	rule := HeaderRule{}
	// Expect it to panic (GetHeaderValue not implemented)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	rule.GetHeaderValue()
}

func TestHeaderRuleValidateNilSource(t *testing.T) {
	rule := HeaderRule{
		destinationHeaders: &http.Header{},
	}
	_, err := rule.Validate()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "invalid source header" {
		t.Errorf("Expected error 'invalid source header', got: %v", err)
	}
}

func TestHeaderRuleValidateNilDestination(t *testing.T) {
	rule := HeaderRule{
		sourceHeaders: &http.Header{},
	}
	_, err := rule.Validate()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "invalid destination header" {
		t.Errorf("Expected error 'invalid destination header', got: %v", err)
	}
}

func TestCopyHeaderRuleSetValue(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	req.Header.Set("Test-Header", "testValue")
	rule, _ := NewCopyHeaderRule("Test-Header", "testValue", "request", "response", req, res)

	rule.SetHeader()
	// Check that the value on the header has been correctly set
	if res.Header.Get(rule.Header) != req.Header.Get(rule.Header) {
		t.Errorf("Expected header value to be '%s', got: %s", req.Header.Get(rule.Header), res.Header.Get(rule.Header))
	}
}

