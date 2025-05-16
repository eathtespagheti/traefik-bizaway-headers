package headerrules

import (
	"bizaway/headerrules/internal/config"
	"bizaway/headerrules/internal/rules"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"
)

type HeadersPlugin struct {
	next                         http.Handler
	requestHeadersConfiguration  map[string]string
	responseHeadersConfiguration map[string]string
	name                         string
	template                     *template.Template
	// Array of headers rules to apply
	requestHeadersRules    []rules.HeaderRule
	responseHeadersRules   []rules.HeaderRule
}

// Take in a configuraiton map for a set of header rules and creates an array of said headerrules
// func createHeaderRule(configuration map[string]string, request *http.Request, response *http.Response, destination string) rules.HeaderRuleInterface {
	
// }

// New created a new Headers plugin.
func New(ctx context.Context, next http.Handler, config *config.Config, name string) (http.Handler, error) {
	if len(config.Headers.Request) == 0 && len(config.Headers.Response) == 0 {
		return nil, fmt.Errorf("at least one request or response header must be configured")
	}	

	// Initialize the headers plugin
	return &HeadersPlugin{
		requestHeadersConfiguration:  config.Headers.Request,
		responseHeadersConfiguration: config.Headers.Response,
		next:                         next,
		name:                         name,
		template:                     template.New("header-value").Delims("[[", "]]"),
	}, nil
}

func (a *HeadersPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Modify Request Headers
	for key, value := range a.requestHeadersConfiguration {
		tmpl, err := a.template.Parse(value)
		if err != nil {
			http.Error(rw, fmt.Sprintf("error parsing request header template '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		writer := &bytes.Buffer{}
		err = tmpl.Execute(writer, req)
		if err != nil {
			http.Error(rw, fmt.Sprintf("error executing request header template '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		req.Header.Set(key, writer.String())
	}

	// If the Host header has been set, copy it over the req Host
	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	// Wrap the ResponseWriter to capture headers set later.
	wrappedWriter := &responseWriterWrapper{
		ResponseWriter:  rw,
		responseHeaders: a.responseHeadersConfiguration,
		template:        a.template,
		originalHeaders: rw.Header(),
		request:         req,
	}
	// Pass the wrapped writer to the next handler.
	a.next.ServeHTTP(wrappedWriter, req)

	//response header will be set inside the responseWriterWrapper
}

// responseWriterWrapper is a custom http.ResponseWriter that captures response headers.
type responseWriterWrapper struct {
	http.ResponseWriter
	responseHeaders map[string]string
	template        *template.Template
	originalHeaders http.Header
	request         *http.Request
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	// Create the aggregated data for the template interpreter
	var templateData = struct {
		Headers *http.Header
		Request *http.Request
	}{
		Headers: &w.originalHeaders,
		Request: w.request,
	}

	for key, value := range w.responseHeaders {
		tmpl, err := w.template.Parse(value)
		if err != nil {
			http.Error(w.ResponseWriter, fmt.Sprintf("error parsing response header template '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		writer := &bytes.Buffer{}
		err = tmpl.Execute(writer, templateData)
		if err != nil {
			http.Error(w.ResponseWriter, fmt.Sprintf("error executing response header template '%s': %v", key, err), http.StatusInternalServerError)
			return
		}
		w.ResponseWriter.Header().Set(key, writer.String())
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
