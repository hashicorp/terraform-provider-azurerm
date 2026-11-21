// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func TestExpandKeyVaultCertificateImportParameters(t *testing.T) {
	certificate := KeyVaultCertificateImportParameters{
		CertificateData:     "aGVsbG8=",
		CertificatePassword: "password",
	}
	policy := &keyvault.CertificatePolicy{}

	testData := []struct {
		name     string
		tags     map[string]any
		policy   *keyvault.CertificatePolicy
		expected keyvault.CertificateImportParameters
	}{
		{
			name: "with tags",
			tags: map[string]any{"environment": "test"},
			expected: keyvault.CertificateImportParameters{
				Base64EncodedCertificate: pointer.To("aGVsbG8="),
				Password:                 pointer.To("password"),
				Tags:                     map[string]*string{"environment": pointer.To("test")},
			},
		},
		{
			name: "without tags",
			tags: map[string]any{},
			expected: keyvault.CertificateImportParameters{
				Base64EncodedCertificate: pointer.To("aGVsbG8="),
				Password:                 pointer.To("password"),
				Tags:                     map[string]*string{},
			},
		},
		{
			name: "with nil tags",
			tags: nil,
			expected: keyvault.CertificateImportParameters{
				Base64EncodedCertificate: pointer.To("aGVsbG8="),
				Password:                 pointer.To("password"),
				Tags:                     map[string]*string{},
			},
		},
		{
			name:   "with a certificate policy",
			tags:   map[string]any{},
			policy: policy,
			expected: keyvault.CertificateImportParameters{
				Base64EncodedCertificate: pointer.To("aGVsbG8="),
				Password:                 pointer.To("password"),
				CertificatePolicy:        policy,
				Tags:                     map[string]*string{},
			},
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			actual := expandKeyVaultCertificateImportParameters(certificate, v.policy, v.tags)
			if diff := cmp.Diff(v.expected, actual); diff != "" {
				t.Errorf("unexpected diff (-expected +actual):\n%s", diff)
			}
		})
	}
}
