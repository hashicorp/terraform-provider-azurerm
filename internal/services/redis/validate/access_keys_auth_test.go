package validate

import (
	"testing"
)

func TestValidateAccessKeysAuth_valid(t *testing.T) {
	if err := ValidateAccessKeysAuth(true, true); err != nil {
		t.Fatalf("Should be valid if accessKeysAuthenticationDisabled: true and activeDirectoryAuthenticationEnabled: true but got error: %v", err)
	}
}

func TestValidateAccessKeysAuth_invalid(t *testing.T) {
	if err := ValidateAccessKeysAuth(true, false); err == nil {
		t.Fatalf("Should return error if accessKeysAuthenticationDisabled: true and activeDirectoryAuthenticationEnabled: false")
	}
}
