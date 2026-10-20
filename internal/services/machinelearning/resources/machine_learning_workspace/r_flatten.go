// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
)

func flattenMachineLearningWorkspaceIdentity(input *identity.LegacySystemAndUserAssignedMap) (*[]any, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        input.Type,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		// api uses `SystemAssigned,UserAssigned` (no space), so normalize it back
		if input.Type == "SystemAssigned,UserAssigned" {
			transform.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}

		if input.IdentityIds != nil {
			for k, v := range input.IdentityIds {
				transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientId,
					PrincipalId: v.PrincipalId,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenMachineLearningWorkspaceEncryption(input *workspaces.EncryptionProperty) (*[]any, error) {
	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return &[]any{}, nil
	}

	keyVaultId := ""
	keyVaultKeyId := ""

	if input.KeyVaultProperties.KeyIdentifier != "" {
		keyVaultKeyId = input.KeyVaultProperties.KeyIdentifier
	}
	if input.KeyVaultProperties.KeyVaultArmId != "" {
		keyVaultId = input.KeyVaultProperties.KeyVaultArmId
	}

	userAssignedIdentityId := ""
	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, fmt.Errorf("parsing userAssignedIdentityId %q: %+v", *input.Identity.UserAssignedIdentity, err)
		}

		userAssignedIdentityId = id.ID()
	}

	return &[]any{
		map[string]any{
			"user_assigned_identity_id": userAssignedIdentityId,
			"key_vault_id":              keyVaultId,
			"key_id":                    keyVaultKeyId,
		},
	}, nil
}

func flattenMachineLearningWorkspaceFeatureStore(input *workspaces.FeatureStoreSettings) *[]any {
	if input == nil {
		return &[]any{}
	}

	computerSparkRunTimeVersion := ""
	offlineConnectionName := ""
	onlineConnectionName := ""

	if input.ComputeRuntime != nil && input.ComputeRuntime.SparkRuntimeVersion != nil {
		computerSparkRunTimeVersion = *input.ComputeRuntime.SparkRuntimeVersion
	}
	if input.OfflineStoreConnectionName != nil {
		offlineConnectionName = *input.OfflineStoreConnectionName
	}

	if input.OnlineStoreConnectionName != nil {
		onlineConnectionName = *input.OnlineStoreConnectionName
	}

	return &[]any{
		map[string]any{
			"computer_spark_runtime_version": computerSparkRunTimeVersion,
			"offline_connection_name":        offlineConnectionName,
			"online_connection_name":         onlineConnectionName,
		},
	}
}

func flattenMachineLearningWorkspaceManagedNetwork(i *workspaces.ManagedNetworkSettings, provisionNetworkNow *bool) *[]any {
	if i == nil {
		return &[]any{}
	}

	out := map[string]any{}

	if i.IsolationMode != nil {
		out["isolation_mode"] = *i.IsolationMode
	}
	out["provision_on_creation_enabled"] = pointer.From(provisionNetworkNow)

	return &[]any{out}
}

func flattenMachineLearningWorkspaceServerlessCompute(i *workspaces.ServerlessComputeSettings) *[]any {
	if i == nil {
		return &[]any{}
	}

	out := map[string]any{}

	if i.ServerlessComputeCustomSubnet != nil {
		out["subnet_id"] = *i.ServerlessComputeCustomSubnet
	}

	if i.ServerlessComputeNoPublicIP != nil {
		out["public_ip_enabled"] = !*i.ServerlessComputeNoPublicIP
	}

	return &[]any{out}
}
