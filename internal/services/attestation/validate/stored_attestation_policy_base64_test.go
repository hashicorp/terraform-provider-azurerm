// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestContainsABase64UriEncodedJWTOfAStoredAttestationPolicy(t *testing.T) {
	testData := []struct {
		expectError bool
		input       string
	}{
		{
			expectError: true,
			input:       "",
		},
		{
			expectError: true,
			input:       "1.2.3",
		},
		{
			expectError: true,
			input:       "{}.{}.{}",
		},
		{
			expectError: true,
			input:       base64.StdEncoding.EncodeToString([]byte("{}.{}.{}")),
		},
		{
			expectError: true,
			input:       base64.URLEncoding.EncodeToString([]byte("{}.{}.{}")),
		},
		{
			expectError: true,
			input:       base64.RawURLEncoding.EncodeToString([]byte("{}.{}.{}")),
		},
		{
			expectError: true,
			input: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"ohai": "there",
				})
				str, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return str
			}(),
		},
		{
			expectError: false,
			input:       generateJWT("bob"),
		},
	}
	for i, item := range testData {
		t.Logf("Testing Item %d..", i)

		warnings, errors := ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy(item.input, "some_field")
		if len(warnings) != 0 {
			t.Fatalf("should have no warnings but got %d: %+v", len(warnings), warnings)
		}
		if item.expectError {
			if len(errors) > 0 {
				continue
			}

			t.Fatalf("expected an error but didn't get one")
		}

		if len(errors) != 0 {
			t.Fatalf("expected no error but got %d: %+v", len(errors), errors)
		}
	}
}

func generateJWT(name string) string {
	// document about create policy: https://learn.microsoft.com/en-us/azure/attestation/author-sign-policy
	policyContent := `version=1.0;
authorizationrules
{
[type=="secureBootEnabled", value==true, issuer=="AttestationService"]=>permit();
};

issuancerules
{
=> issue(type="SecurityLevelValue", value=100);
};`
	b64Encoded := base64.RawURLEncoding.EncodeToString([]byte(policyContent))
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"AttestationPolicy": b64Encoded,
	})
	token.Header["jku"] = fmt.Sprintf("https://%s.uks.attest.azure.net/certs", name)
	token.Header["kid"] = "xxx"

	str, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return str
}
