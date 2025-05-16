package rules

import (
	"net/http"
	"text/template"
)

// Generic HeaderRule interface meant to be extended into specific HeaderRules
type HeaderRuleInterface interface {
	// Get the value that should be set in the Header
	GetHeaderValue() string
	// Set the Header value
	SetHeader()
	// Validate that all the Header Rule parameters are correctly configured, if something is wrong return an error
	Validate() (bool, error)

	// Initialize the sourceHeader array
	initSourceHeader()
	// Initialize the destinationHeader array
	initDestinationHeader()
	// Initialize the template
	initTemplate()
	// Function that determines wether or not the header modification should happen, by default always true
	headerChangeConditions() bool
	// Parse the value field as a template
	parseValueTemplate() (string, error)
}

// Generic HeaderRule meant to be extended into specific HeaderRules
type HeaderRule struct {
	HeaderRuleInterface
	Request            *http.Request
	Response           *http.Response
	Header             string
	Value              string
	Source             string
	Destination        string
	ChangeCondition    func(*HeaderRule) bool
	EnableTemplates    bool
	sourceHeaders      *http.Header
	destinationHeaders *http.Header
	template           *template.Template
}

// String Rule
type StringHeaderRule struct {
	HeaderRule
}

// Copy Rule
type CopyHeaderRule struct {
	HeaderRule
}
