package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GIVEN_BasicAuthSet_WHEN_CallingIsAllowedWithCorrectCreds_THEN_TrueReturned(t *testing.T) {
	// assemble
	username := "u1"
	password := "p1"
	request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	request.SetBasicAuth(username, password)
	provider := NewBasicAuthProvider(username, password)

	// act
	result := provider.IsAllowed(request)

	// assert
	if !result {
		t.Errorf("expected request to be allowed, but failed")
	}
}

func Test_GIVEN_BasicAuthSet_WHEN_CallingIsAllowedWithoutCreds_THEN_FalseReturned(t *testing.T) {
	// assemble
	request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	provider := NewBasicAuthProvider("u1", "p1")

	// act
	result := provider.IsAllowed(request)

	// assert
	if result {
		t.Errorf("expected request to be denied, but was allowed")
	}
}

func Test_GIVEN_BasicAuthSet_WHEN_CallingIsAllowedWithWrongCreds_THEN_FalseReturned(t *testing.T) {
	// assemble
	request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	request.SetBasicAuth("wrong", "pw")
	provider := NewBasicAuthProvider("u1", "p1")

	// act
	result := provider.IsAllowed(request)

	// assert
	if result {
		t.Errorf("expected request to be denied, but was allowed")
	}
}
