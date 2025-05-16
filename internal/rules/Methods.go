package rules

import (
	"bytes"
	"errors"
	"net/http"
	"text/template"
)

// HeaderRule

// Function that determines wether or not the header modification should happen, by default always true
func (h *HeaderRule) headerChangeConditions() bool {
	// Check if the change condition function is present
	if h.ChangeCondition != nil {
		return h.ChangeCondition(h)
	}
	return true
}

// Default GetHeaderValue that panics since an override method is not defined
func (h *HeaderRule) GetHeaderValue() string {
	panic("GetHeaderValue not implemented")
}

// Set the Header value
func (h *HeaderRule) SetHeader() {
	if !h.headerChangeConditions() {
		return
	}
	h.destinationHeaders.Set(h.Header, h.GetHeaderValue())
}

// Initialize the sourceHeader array
func (h *HeaderRule) initSourceHeader() {
	if h.Source == "request" {
		h.sourceHeaders = &h.Request.Header
	} else if h.Source == "response" {
		h.sourceHeaders = &h.Response.Header
	}
}

// Initialize the destinationHeader array
func (h *HeaderRule) initDestinationHeader() {
	if h.Destination == "request" {
		h.destinationHeaders = &h.Request.Header
	} else if h.Destination == "response" {
		h.destinationHeaders = &h.Response.Header
	}
}

// Initialize the template
func (h *HeaderRule) initTemplate() {
	h.template = template.New("header-value").Delims("[[", "]]")
}

// parseValueTemplate
func (h *HeaderRule) parseValueTemplate() (string, error) {
	tmpl, err := h.template.Parse(h.Value)
	if err != nil {
		return "", err
	}

	writer := &bytes.Buffer{}
	err = tmpl.Execute(writer, h)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

// Validate that all the Header Rule parameters are correctly configured
func (h *HeaderRule) Validate() (bool, error) {
	// Check source header
	if h.sourceHeaders == nil {
		return false, errors.New("invalid source header")
	}

	// Check destination header
	if h.destinationHeaders == nil {
		return false, errors.New("invalid destination header")
	}

	return true, nil
}

// String Rule

// Get the the Header name
func (s *StringHeaderRule) GetHeaderValue() string {
	if s.EnableTemplates {
		value, err := s.parseValueTemplate()
		if err != nil {
			return s.Value
		}
		return value
	}
	return s.Value
}

// Set the Header value
func (h *StringHeaderRule) SetHeader() {
	if !h.headerChangeConditions() {
		return
	}
	h.destinationHeaders.Set(h.Header, h.GetHeaderValue())
}

// Create a new StringHeaderRule
func NewStringHeaderRule(header string, value string, destination string, request *http.Request, response *http.Response, headerChangeCondition func(*HeaderRule) bool, enableTemplates bool) (*StringHeaderRule, error) {
	var shr *StringHeaderRule = &StringHeaderRule{
		HeaderRule: HeaderRule{
			Header:          header,
			Value:           value,
			Destination:     destination,
			Request:         request,
			Response:        response,
			ChangeCondition: headerChangeCondition,
			EnableTemplates: enableTemplates,
		},
	}

	// Init the destinationHeader array
	shr.initDestinationHeader()

	// Set the source header the same as the destination header in order to pass validation
	shr.sourceHeaders = shr.destinationHeaders

	// Init the template if enabled
	if enableTemplates {
		shr.initTemplate()
	}

	// Validate the rule
	_, err := shr.Validate()
	return shr, err
}

// Copy Rule

// Get the the Header name
func (h *CopyHeaderRule) GetHeaderValue() string {
	return h.sourceHeaders.Get(h.Header)
}

// Set the Header value
func (h *CopyHeaderRule) SetHeader() {
	if !h.headerChangeConditions() {
		return
	}
	h.destinationHeaders.Set(h.Header, h.GetHeaderValue())
}

// Create a new CopyHeaderRule
func NewCopyHeaderRule(header string, value string, source string, destination string, request *http.Request, response *http.Response) (*CopyHeaderRule, error) {
	var chr *CopyHeaderRule = &CopyHeaderRule{
		HeaderRule: HeaderRule{
			Header:      header,
			Value:       value,
			Source:      source,
			Destination: destination,
			Request:     request,
			Response:    response,
		},
	}

	// Init the Header arrays
	chr.initSourceHeader()
	chr.initDestinationHeader()

	// Validate the rule
	_, err := chr.Validate()
	return chr, err
}
