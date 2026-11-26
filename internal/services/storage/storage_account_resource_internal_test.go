// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
)

func TestNormalizeAzureFilesAuthentication(t *testing.T) {
	testData := []struct {
		name     string
		input    *storageaccounts.AzureFilesIdentityBasedAuthentication
		expected *storageaccounts.AzureFilesIdentityBasedAuthentication
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "no active directory properties",
			input: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsAADDS,
			},
			expected: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsAADDS,
			},
		},
		{
			name: "aadds with an empty domain guid",
			input: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsAADDS,
				ActiveDirectoryProperties: &storageaccounts.ActiveDirectoryProperties{
					DomainGuid: pointer.To(""),
					DomainName: pointer.To("contoso.com"),
				},
			},
			expected: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsAADDS,
				ActiveDirectoryProperties: &storageaccounts.ActiveDirectoryProperties{
					DomainName: pointer.To("contoso.com"),
				},
			},
		},
		{
			name: "domain guid set",
			input: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				ActiveDirectoryProperties: &storageaccounts.ActiveDirectoryProperties{
					DomainGuid: pointer.To("12345678-1234-1234-1234-123456789012"),
					DomainName: pointer.To("contoso.com"),
				},
			},
			expected: &storageaccounts.AzureFilesIdentityBasedAuthentication{
				ActiveDirectoryProperties: &storageaccounts.ActiveDirectoryProperties{
					DomainGuid: pointer.To("12345678-1234-1234-1234-123456789012"),
					DomainName: pointer.To("contoso.com"),
				},
			},
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			actual := normalizeAzureFilesAuthentication(v.input)
			if diff := cmp.Diff(v.expected, actual); diff != "" {
				t.Errorf("unexpected diff (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestExpandAccountActiveDirectoryPropertiesWithAnEmptyDomainGuid(t *testing.T) {
	actual := expandAccountActiveDirectoryProperties([]any{
		map[string]any{
			"domain_guid":         "",
			"domain_name":         "contoso.com",
			"storage_sid":         "",
			"domain_sid":          "",
			"forest_name":         "",
			"netbios_domain_name": "",
		},
	})

	if actual == nil {
		t.Fatal("expected a value, got nil")
	}
	if actual.DomainGuid != nil {
		t.Fatalf("expected `DomainGuid` to be nil, got %q", pointer.From(actual.DomainGuid))
	}
	if pointer.From(actual.DomainName) != "contoso.com" {
		t.Fatalf("expected `DomainName` to be `contoso.com`, got %q", pointer.From(actual.DomainName))
	}
}
