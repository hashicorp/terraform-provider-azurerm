// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestKustoClusterCustomerManagedKeyExpandFlatten(t *testing.T) {
	input := []any{
		map[string]any{
			"key_vault_key_id":             "https://example.vault.azure.net/keys/example-key/version",
			"user_assigned_identity_id":    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example",
			"federated_identity_client_id": "00000000-0000-0000-0000-000000000000",
		},
	}

	expanded, err := expandKustoClusterCustomerManagedKey(input)
	if err != nil {
		t.Fatalf("expanding customer managed key: %s", err)
	}
	if got, want := pointer.From(expanded.FederatedIdentityClientId), "00000000-0000-0000-0000-000000000000"; got != want {
		t.Fatalf("federated identity client ID was %q, want %q", got, want)
	}

	flattened, err := flattenKustoClusterCustomerManagedKey(expanded)
	if err != nil {
		t.Fatalf("flattening customer managed key: %s", err)
	}
	if !reflect.DeepEqual(flattened, input) {
		t.Fatalf("flattened customer managed key was %#v, want %#v", flattened, input)
	}
}
