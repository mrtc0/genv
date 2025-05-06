package mocks

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	MockStsGetCallerIdentityValidEndpoint = &MockEndpoint{
		Request: &MockRequest{
			Body: url.Values{
				"Action":  []string{"GetCallerIdentity"},
				"Version": []string{"2011-06-15"},
			}.Encode(),
			Method: http.MethodPost,
			Uri:    "/",
		},
		Response: &MockResponse{
			Body:        MockStsGetCallerIdentityValidResponseBody,
			ContentType: "text/xml",
			StatusCode:  http.StatusOK,
		},
	}

	MockStsGetRoleCredentialsValidEndpoint = &MockEndpoint{
		Request: &MockRequest{
			Method: http.MethodGet,
			Uri: fmt.Sprintf(
				"/federation/credentials?account_id=%s&role_name=%s",
				StsGetRoleCredentialsAccountId,
				StsGetRoleCredentialsRoleName,
			),
		},
		Response: &MockResponse{
			Body: fmt.Sprintf(
				MockStsGetRoleCredentialsValidResponseBodyTemplate,
				StsGetRoleCredentialsAccessKeyId,
				time.Now().Add(15*time.Minute).UnixNano()/int64(time.Millisecond),
				StsGetRoleCredentialsSecretAccessKey,
				StsGetRoleCredentialsSessionToken,
			),
			ContentType: "application/json",
			StatusCode:  http.StatusOK,
		},
	}
)

// MockEndpoint represents a mock endpoint for testing purposes.
type MockEndpoint struct {
	Request  *MockRequest
	Response *MockResponse
}

// MockRequest represents a mock HTTP request for testing purposes.
type MockRequest struct {
	Method string
	Uri    string
	Body   string
}

// MockResponse represents a mock HTTP response for testing purposes.
type MockResponse struct {
	StatusCode  int
	Body        string
	ContentType string
}

func MockAwsApiServer(t *testing.T, endpoints []*MockEndpoint) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		requestBody := buf.String()

		for _, endpoint := range endpoints {
			if r.Method == endpoint.Request.Method && r.RequestURI == endpoint.Request.Uri && requestBody == endpoint.Request.Body {
				w.Header().Set("Content-Type", endpoint.Response.ContentType)
				w.Header().Set("X-Amzn-Requestid", "mock-request-id")
				w.Header().Set("Date", time.Now().Format(time.RFC1123))
				w.WriteHeader(endpoint.Response.StatusCode)

				fmt.Fprintln(w, endpoint.Response.Body)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		t.Logf("No matching endpoint found: %s %s %s", r.Method, r.RequestURI, requestBody)
	}))

	return server
}
