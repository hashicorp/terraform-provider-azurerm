// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s : %+v", id.String(), err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
			state := SiteRecoveryReplicationPolicyVmwareAssociationModel{
				Name:            id.ReplicationProtectionContainerMappingName,
				RecoveryVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					policyId := ""
					// tracked on https://github.com/Azure/azure-rest-api-specs/issues/24751
					if prop.PolicyId != nil && *prop.PolicyId != "" {
						parsedPolicyId, err := replicationpolicies.ParseReplicationPolicyIDInsensitively(*prop.PolicyId)
						if err != nil {
							return fmt.Errorf("parsing %s: %+v", *prop.PolicyId, err)
						}
						policyId = parsedPolicyId.ID()
					}

					state.RecoveryReplicationPolicyId = policyId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const SiteRecoveryReplicationPolicyVMWareAssociationTargetContainerId string = "Microsoft Azure"

func validateReplicationPolicyAssociationVaultConfig(ctx context.Context, settingsClient *replicationvaultsetting.ReplicationVaultSettingClient, vaultId string) error {
	vId, err := replicationvaultsetting.ParseVaultID(vaultId)
	if err != nil {
		return fmt.Errorf("parse %s: %+v", vaultId, err)
	}

	settingsId := replicationvaultsetting.NewReplicationVaultSettingID(vId.SubscriptionId, vId.ResourceGroupName, vId.VaultName, "default")
	resp, err := settingsClient.Get(ctx, settingsId)
	if err != nil {
		return fmt.Errorf("retire %s: %+v", settingsId, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.VMwareToAzureProviderType != nil {
		return fmt.Errorf("can not associate a modern policy to classic VMWare vault")
	}

	return nil
}

func fetchReplicationPolicyAssociationContainerId(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, vaultId string) (containerId string, err error) {
	vId, err := replicationprotectioncontainers.ParseVaultID(vaultId)
	if err != nil {
		return "", fmt.Errorf("parse %s: %+v", vaultId, err)
	}

	resp, err := containerClient.ListComplete(ctx, *vId)
	if err != nil {
		return "", err
	}

	if len(resp.Items) != 1 {
		return "", fmt.Errorf("there should be only one protection container in Classic Recovery Vault, get: %v", len(resp.Items))
	}

	return helpers.HandleAzureSdkForGoBug2824(*resp.Items[0].Id), nil
}
