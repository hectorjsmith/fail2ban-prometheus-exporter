package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testAuthProvider struct {
	match bool
}

func (p testAuthProvider) IsAllowed(request *http.Request) bool {
	return p.match
}

func newTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "http://example.com", nil)
}

func executeAuthMiddlewareTest(t *testing.T, authMatches bool, expectedCode int, expectedCallCount int) {
	callCount := 0
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
	}

	handler := AuthMiddleware(testHandler, testAuthProvider{match: authMatches})
	recorder := httptest.NewRecorder()
	request := newTestRequest()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != expectedCode {
		t.Errorf("statusCode = %v, want %v", recorder.Code, expectedCode)
	}
	if callCount != expectedCallCount {
		t.Errorf("callCount = %v, want %v", callCount, expectedCallCount)
	}
}

func Test_GIVEN_MatchingBasicAuth_WHEN_MethodCalled_THEN_RequestProcessed(t *testing.T) {
	executeAuthMiddlewareTest(t, true, http.StatusOK, 1)
}

func Test_GIVEN_NonMatchingBasicAuth_WHEN_MethodCalled_THEN_RequestRejected(t *testing.T) {
	executeAuthMiddlewareTest(t, false, http.StatusUnauthorized, 0)
}
