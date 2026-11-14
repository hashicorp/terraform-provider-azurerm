// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
)

func TestKustoClusterCustomerManagedKeyExpandFlatten(t *testing.T) {
	input := []any{
		map[string]any{
			"key_vault_key_id":             "https://example.vault.azure.net/keys/example-key/version",
			"user_assigned_identity_id":    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example",
			"federated_identity_client_id": "00000000-0000-0000-0000-000000000000",
		},
	}

	expanded, federatedIdentityClientID, err := expandKustoClusterCustomerManagedKey(input)
	if err != nil {
		t.Fatalf("expanding customer managed key: %s", err)
	}
	if got, want := pointer.From(federatedIdentityClientID), "00000000-0000-0000-0000-000000000000"; got != want {
		t.Fatalf("federated identity client ID was %q, want %q", got, want)
	}

	flattened, err := flattenKustoClusterCustomerManagedKey(expanded, federatedIdentityClientID)
	if err != nil {
		t.Fatalf("flattening customer managed key: %s", err)
	}
	if !reflect.DeepEqual(flattened, input) {
		t.Fatalf("flattened customer managed key was %#v, want %#v", flattened, input)
	}
}

func TestMarshalKustoClusterRequestAddsFederatedIdentityClientID(t *testing.T) {
	input := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{
				KeyName: pointer.To("example-key"),
			},
		},
	}
	clientID := "00000000-0000-0000-0000-000000000000"

	body, err := marshalKustoClusterRequest(input, &clientID)
	if err != nil {
		t.Fatalf("marshalling Kusto cluster request: %s", err)
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("unmarshalling Kusto cluster request: %s", err)
	}
	properties := request["properties"].(map[string]any)
	keyVaultProperties := properties["keyVaultProperties"].(map[string]any)
	if got := keyVaultProperties["federatedIdentityClientId"]; got != clientID {
		t.Fatalf("federatedIdentityClientId was %#v, want %q", got, clientID)
	}
}

func TestMarshalKustoClusterRequestClearsCustomerManagedKey(t *testing.T) {
	input := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{},
		},
	}

	body, err := marshalKustoClusterRequest(input, nil)
	if err != nil {
		t.Fatalf("marshalling Kusto cluster request: %s", err)
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("unmarshalling Kusto cluster request: %s", err)
	}
	properties := request["properties"].(map[string]any)
	keyVaultProperties := properties["keyVaultProperties"].(map[string]any)
	if len(keyVaultProperties) != 0 {
		t.Fatalf("keyVaultProperties was %#v, want an empty object to clear the key", keyVaultProperties)
	}
}

func TestUpdateKustoClusterRestoresAutoscaleAfterFailure(t *testing.T) {
	skuErr := errors.New("SKU update failed")
	var restored bool

	err := updateKustoClusterAndRestoreAutoscaleOnFailure(
		func() error { return skuErr },
		func() error {
			restored = true
			return nil
		},
	)
	if err != skuErr {
		t.Fatalf("error was %v, want %v", err, skuErr)
	}
	if !restored {
		t.Fatal("expected autoscale to be restored after SKU update failure")
	}
}

func TestUpdateKustoClusterReportsAutoscaleRestoreFailure(t *testing.T) {
	skuErr := errors.New("SKU update failed")
	restoreErr := errors.New("autoscale restore failed")

	err := updateKustoClusterAndRestoreAutoscaleOnFailure(
		func() error { return skuErr },
		func() error { return restoreErr },
	)
	if !errors.Is(err, skuErr) {
		t.Fatalf("error %v does not contain SKU update failure", err)
	}
	if !errors.Is(err, restoreErr) {
		t.Fatalf("error %v does not contain autoscale restore failure", err)
	}
}

func TestKustoClusterCustomerManagedKeyCanBeRemoved(t *testing.T) {
	schema := kustoClusterSchema()["customer_managed_key"]
	if !schema.Optional || schema.Computed {
		t.Fatal("customer_managed_key must be Optional and not Computed to allow explicit removal")
	}
}
