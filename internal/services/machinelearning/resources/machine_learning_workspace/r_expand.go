// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
)

func expandMachineLearningWorkspaceIdentity(input []any) (*identity.LegacySystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.LegacySystemAndUserAssignedMap{
		Type: expanded.Type,
	}
	if len(expanded.IdentityIds) > 0 {
		out.IdentityIds = map[string]identity.UserAssignedIdentityDetails{}
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func expandMachineLearningWorkspaceEncryption(input []any) *workspaces.EncryptionProperty {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: raw["key_vault_id"].(string),
			KeyIdentifier: raw["key_id"].(string),
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if raw["user_assigned_identity_id"].(string) != "" {
		out.Identity.UserAssignedIdentity = pointer.To(raw["user_assigned_identity_id"].(string))
	}

	return &out
}

func expandMachineLearningWorkspaceFeatureStore(input []any) *workspaces.FeatureStoreSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	out := workspaces.FeatureStoreSettings{}

	if raw["computer_spark_runtime_version"].(string) != "" {
		out.ComputeRuntime = &workspaces.ComputeRuntimeDto{
			SparkRuntimeVersion: pointer.To(raw["computer_spark_runtime_version"].(string)),
		}
	}

	if raw["offline_connection_name"].(string) != "" {
		out.OfflineStoreConnectionName = pointer.To(raw["offline_connection_name"].(string))
	}

	if raw["online_connection_name"].(string) != "" {
		out.OnlineStoreConnectionName = pointer.To(raw["online_connection_name"].(string))
	}
	return &out
}

func expandMachineLearningWorkspaceManagedNetwork(i []any) (*workspaces.ManagedNetworkSettings, bool) {
	if len(i) == 0 || i[0] == nil {
		return nil, false
	}

	v := i[0].(map[string]any)

	return &workspaces.ManagedNetworkSettings{
		IsolationMode: pointer.ToEnum[workspaces.IsolationMode](v["isolation_mode"].(string)),
	}, v["provision_on_creation_enabled"].(bool)
}

func expandMachineLearningWorkspaceServerlessCompute(i []any) *workspaces.ServerlessComputeSettings {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	v := i[0].(map[string]any)

	serverlessCompute := workspaces.ServerlessComputeSettings{
		ServerlessComputeNoPublicIP: pointer.To(!v["public_ip_enabled"].(bool)),
	}

	if subnetId, ok := v["subnet_id"].(string); ok && subnetId != "" {
		serverlessCompute.ServerlessComputeCustomSubnet = pointer.To(subnetId)
	}

	return &serverlessCompute
}
