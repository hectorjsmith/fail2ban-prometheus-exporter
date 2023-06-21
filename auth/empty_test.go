package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GIVEN_EmptyAuth_WHEN_CallingIsAllowedWithoutAuth_THEN_TrueReturned(t *testing.T) {
	// assemble
	request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	provider := NewEmptyAuthProvider()

	// act
	response := provider.IsAllowed(request)

	// assert
	if !response {
		t.Errorf("expected request to be allowed, but failed")
	}
}

func Test_GIVEN_EmptyAuth_WHEN_CallingIsAllowedWithAuth_THEN_TrueReturned(t *testing.T) {
	// assemble
	request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	request.SetBasicAuth("user", "pass")
	provider := NewEmptyAuthProvider()

	// act
	response := provider.IsAllowed(request)

	// assert
	if !response {
		t.Errorf("expected request to be allowed, but failed")
	}
}
