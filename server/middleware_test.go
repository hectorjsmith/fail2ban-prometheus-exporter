package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testAuthProvider struct {
	enabled bool
	match   bool
}

func (p testAuthProvider) Enabled() bool {
	return p.enabled
}

func (p testAuthProvider) DoesBasicAuthMatch(username, password string) bool {
	return p.match
}

func newTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "http://example.com", nil)
}

func executeBasicAuthMiddlewareTest(t *testing.T, authEnabled bool, authMatches bool, expectedCode int, expectedCallCount int) {
	callCount := 0
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
	}

	handler := BasicAuthMiddleware(testHandler, testAuthProvider{enabled: authEnabled, match: authMatches})
	recorder := httptest.NewRecorder()
	request := newTestRequest()
	if authEnabled {
		request.SetBasicAuth("test", "test")
	}
	handler.ServeHTTP(recorder, request)

	if recorder.Code != expectedCode {
		t.Errorf("statusCode = %v, want %v", recorder.Code, expectedCode)
	}
	if callCount != expectedCallCount {
		t.Errorf("callCount = %v, want %v", callCount, expectedCallCount)
	}
}

func Test_GIVEN_DisabledBasicAuth_WHEN_MethodCalled_THEN_RequestProcessed(t *testing.T) {
	executeBasicAuthMiddlewareTest(t, false, false, http.StatusOK, 1)
}

func Test_GIVEN_EnabledBasicAuth_WHEN_MethodCalledWithCorrectCredentials_THEN_RequestProcessed(t *testing.T) {
	executeBasicAuthMiddlewareTest(t, true, true, http.StatusOK, 1)
}

func Test_GIVEN_EnabledBasicAuth_WHEN_MethodCalledWithIncorrectCredentials_THEN_RequestRejected(t *testing.T) {
	executeBasicAuthMiddlewareTest(t, true, false, http.StatusUnauthorized, 0)
}
