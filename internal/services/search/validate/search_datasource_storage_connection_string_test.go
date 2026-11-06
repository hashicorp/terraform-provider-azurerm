// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestSearchDatasourceStorageConnectionString(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "invalid-connection-string",
			Valid: false,
		},
		{
			Input: "some=random;key=value;pairs=here",
			Valid: false,
		},
		{
			Input: "ResourceId=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sa;",
			Valid: true,
		},
		{
			Input: "DefaultEndpointsProtocol=https;AccountName=mystorageaccount;AccountKey=abc123==;",
			Valid: true,
		},
		{
			Input: "DefaultEndpointsProtocol=https;AccountName=mystorageaccount;",
			Valid: false,
		},
		{
			Input: "AccountName=mystorageaccount;AccountKey=abc123==;",
			Valid: false,
		},
		{
			Input: "BlobEndpoint=https://mystorageaccount.blob.core.windows.net;SharedAccessSignature=sv=2021-06-08&sig=abc123;",
			Valid: true,
		},
		{
			Input: "BlobEndpoint=https://mystorageaccount.blob.core.windows.net;",
			Valid: false,
		},
		{
			Input: "ContainerSharedAccessUri=https://mystorageaccount.blob.core.windows.net/container?sv=2021-06-08&sig=abc123",
			Valid: true,
		},
		{
			Input: "DefaultEndpointsProtocol=https;AccountName=mystorageaccount;AccountKey=abc123==;EndpointSuffix=core.windows.net",
			Valid: true,
		},
	}

	for _, tc := range cases {
		_, errors := SearchDatasourceStorageConnectionString(tc.Input, "connection_string")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("expected %t but got %t for input %q", tc.Valid, valid, tc.Input)
		}
	}
}
