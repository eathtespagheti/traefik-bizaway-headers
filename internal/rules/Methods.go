package rules

import (
	"errors"
	"log"
	"net/http"
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

// Set the Header value
func (h *HeaderRule) SetHeader() {
	if !h.headerChangeConditions() {
		return
	}
	
	h.destinationHeader.Set(h.Value, h.GetHeaderValue())
}

// Initialize the sourceHeader array
func (h *HeaderRule) initSourceHeader() {
	if h.Source == "request" {
		h.sourceHeader = &h.Request.Header
	} else if h.Source == "response" {
		h.sourceHeader = &h.Response.Header
	}
}

// Initialize the destinationHeader array
func (h *HeaderRule) initDestinationHeader() {
	if h.Destination == "request" {
		h.destinationHeader = &h.Request.Header
	} else if h.Destination == "response" {
		h.destinationHeader = &h.Response.Header
	}
}

// Validate that all the Header Rule parameters are correctly configured
func (h *HeaderRule) Validate() (bool, error) {
	// Check source header
	if h.sourceHeader == nil {
		return false, errors.New("invalid source header")
	}

	// Check destination header
	if h.destinationHeader == nil {
		return false, errors.New("invalid destination header")
	}

	return true, nil
}

// String Rule

// Get the the Header name
func (s *StringHeaderRule) GetHeader() string {
	return s.Value
}

// Validate that all the Header Rule parameters are correctly configured
func (h *StringHeaderRule) Validate() (bool, error) {
	if h.destinationHeader == nil {
		return false, errors.New("invalid destination header")
	}

	return true, nil
}

func NewStringHeaderRule(header string, value string, destination string, request *http.Request, response *http.Response, headerChangeCondition func(*HeaderRule) bool) *StringHeaderRule {
	var shr *StringHeaderRule = &StringHeaderRule{
		HeaderRule: HeaderRule{
			Header:      header,
			Value:       value,
			Destination: destination,
			Request:     request,
			Response:    response,
			ChangeCondition: headerChangeCondition,
		},
	}

	// Init the destinationHeader array
	shr.initDestinationHeader()

	// Validate the rule
	valid, err := shr.Validate()
	if !valid {
		log.Fatal(err)
	}

	return shr
}

// Copy Rule

// Get the the Header name
func (h *CopyHeaderRule) GetHeader() string {
	return h.sourceHeader.Get(h.Header)
}

func NewCopyHeaderRule(header string, value string, source string, destination string, request *http.Request, response *http.Response) *CopyHeaderRule {
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
	valid, err := chr.Validate()
	if !valid {
		log.Fatal(err)
	}

	return chr
}
