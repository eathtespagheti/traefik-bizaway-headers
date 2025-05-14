package rules

import (
	"errors"
	"net/http"
	"testing"
)

func testHeaderRuleHeaderChangeConditions(t *testing.T) {
	rule := HeaderRule{}
	if !rule.headerChangeConditions() {
		t.Error("HeaderChangeConditions should return true by default")
	}
}

func testHeaderRuleSetHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule := HeaderRule{
		Request:           req,
		Response:          res,
		Value:             "testValue",
		Destination:       "request",
		destinationHeader: &req.Header,
	}

	rule.SetHeader()
	if req.Header.Get(rule.Value) != "" {
		t.Errorf("Expected header to be empty, got: %s", req.Header.Get(rule.Value))
	}

	rule.Header = "Test-Header"
	rule.SetHeader()
	if req.Header.Get(rule.Header) != "testValue" {
		t.Errorf("Expected header value to be 'testValue', got: %s", req.Header.Get(rule.Header))
	}

	// Test with HeaderChangeConditions returning false
	rule.ChangeCondition = func(*HeaderRule) bool { return false }
	rule.SetHeader()
	if req.Header.Get(rule.Header) != "testValue" {
		t.Errorf("Expected header value to be 'testValue' (no change), got: %s", req.Header.Get(rule.Header))
	}
}

func testHeaderRuleInitSourceHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:  req,
		Response: res,
		Source:   "request",
	}
	rule.initSourceHeader()
	if rule.sourceHeader == nil || rule.sourceHeader != &req.Header {
		t.Error("Expected sourceHeader to be initialized to request header")
	}

	rule.Source = "response"
	rule.initSourceHeader()
	if rule.sourceHeader == nil || rule.sourceHeader != &res.Header {
		t.Error("Expected sourceHeader to be initialized to response header")
	}

	rule.Source = "invalid"
	rule.initSourceHeader()
	if rule.sourceHeader == nil {
		t.Error("Expected sourceHeader to be nil")
	}
}

func testHeaderRuleInitDestinationHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:     req,
		Response:    res,
		Destination: "request",
	}
	rule.initDestinationHeader()
	if rule.destinationHeader == nil || rule.destinationHeader != &req.Header {
		t.Error("Expected destinationHeader to be initialized to request header")
	}

	rule.Destination = "response"
	rule.initDestinationHeader()
	if rule.destinationHeader == nil || rule.destinationHeader != &res.Header {
		t.Error("Expected destinationHeader to be initialized to response header")
	}
	rule.Destination = "invalid"
	rule.initDestinationHeader()
	if rule.destinationHeader == nil {
		t.Error("Expected destinationHeader to be nil")
	}
}

func testHeaderRuleValidate(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}

	rule := HeaderRule{
		Request:           req,
		Response:          res,
		destinationHeader: &req.Header,
		sourceHeader:      &res.Header,
	}
	valid, err := rule.Validate()
	if !valid || err != nil {
		t.Errorf("Expected Validate to return true and nil error, got: %v, %v", valid, err)
	}

	rule.sourceHeader = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid source header" {
		t.Errorf("Expected error message 'invalid source header', got: %v", err)
	}

	rule.sourceHeader = &res.Header
	rule.destinationHeader = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid destination header" {
		t.Errorf("Expected error message 'invalid destination header', got: %v", err)
	}
}

func testStringHeaderRuleGetHeader(t *testing.T) {
	rule := StringHeaderRule{
		HeaderRule: HeaderRule{
			Value: "testHeader",
		},
	}
	if rule.GetHeader() != "testHeader" {
		t.Errorf("Expected GetHeader to return 'testHeader', got: %s", rule.GetHeader())
	}
}

func testStringHeaderRuleValidate(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	rule := StringHeaderRule{
		HeaderRule: HeaderRule{
			destinationHeader: &req.Header,
		},
	}
	valid, err := rule.Validate()
	if !valid || err != nil {
		t.Errorf("Expected Validate to return true and nil error, got: %v, %v", valid, err)
	}

	rule.destinationHeader = nil
	valid, err = rule.Validate()
	if valid || err == nil {
		t.Errorf("Expected Validate to return false and an error, got: %v, %v", valid, err)
	}
	if err != nil && err.Error() != "invalid destination header" {
		t.Errorf("Expected error message 'invalid destination header', got: %v", err)
	}
}

func testNewStringHeaderRule(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule := NewStringHeaderRule("Test-Header", "testValue", "request", req, res, nil)

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
	if rule.destinationHeader != &req.Header {
		t.Error("Expected destinationHeader to be initialized to request header")
	}
}

func testNewStringHeaderRuleInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	NewStringHeaderRule("Test-Header", "testValue", "invalid", req, res, nil)
}

func testCopyHeaderRuleGetHeader(t *testing.T) {
	req := &http.Request{Header: http.Header{"Test-Header": []string{"testValue"}}}
	rule := CopyHeaderRule{
		HeaderRule: HeaderRule{
			Header:       "Test-Header",
			sourceHeader: &req.Header,
		},
	}
	if rule.GetHeader() != "testValue" {
		t.Errorf("Expected GetHeader to return 'testValue', got: %s", rule.GetHeader())
	}
}

func testNewCopyHeaderRule(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	rule := NewCopyHeaderRule("Test-Header", "testValue", "request", "response", req, res)

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
	if rule.sourceHeader != &req.Header {
		t.Error("Expected sourceHeader to be initialized to request header")
	}
	if rule.destinationHeader != &res.Header {
		t.Error("Expected destinationHeader to be initialized to response header")
	}
}

func testNewCopyHeaderRule_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	req := &http.Request{Header: http.Header{}}
	res := &http.Response{Header: http.Header{}}
	NewCopyHeaderRule("Test-Header", "testValue", "invalid", "response", req, res)
}

func testHeaderRuleGetHeaderValue(t *testing.T) {
	rule := HeaderRule{
		Value: "testValue",
	}
	if rule.GetHeaderValue() != "testValue" {
		t.Errorf("Expected GetHeaderValue to return 'testValue', got: %s", rule.GetHeaderValue())
	}
}

func testHeaderRuleGetHeaderValueEmpty(t *testing.T) {
	rule := HeaderRule{}
	if rule.GetHeaderValue() != "" {
		t.Errorf("Expected GetHeaderValue to return '', got: %s", rule.GetHeaderValue())
	}
}

func testHeaderRuleValidateNilSource(t *testing.T) {
	rule := HeaderRule{
		destinationHeader: &http.Header{},
	}
	_, err := rule.Validate()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, errors.New("invalid source header")) {
		t.Errorf("Expected error 'invalid source header', got: %v", err)
	}
}

func testHeaderRuleValidateNilDestination(t *testing.T) {
	rule := HeaderRule{
		sourceHeader: &http.Header{},
	}
	_, err := rule.Validate()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, errors.New("invalid destination header")) {
		t.Errorf("Expected error 'invalid destination header', got: %v", err)
	}
}
