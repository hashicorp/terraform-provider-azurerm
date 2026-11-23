// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SiteRecoveryReplicationPolicyVmwareAssociationModel
			err := metadata.Decode(&model)
			if err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient
			containerClient := metadata.Client.RecoveryServices.ProtectionContainerClient
			settingsClient := metadata.Client.RecoveryServices.VaultsSettingsClient

			if err = validateReplicationPolicyAssociationVaultConfig(ctx, settingsClient, model.RecoveryVaultId); err != nil {
				return fmt.Errorf("validating %s: %+v", model.RecoveryVaultId, err)
			}

			// There should be only 1 fabric and only 1 container in Vmware Vault.
			containerId, err := fetchReplicationPolicyAssociationContainerId(ctx, containerClient, model.RecoveryVaultId)
			if err != nil {
				return fmt.Errorf("fetch replication container from %q: %+v", model.RecoveryVaultId, err)
			}

			parsedContainerId, err := replicationprotectioncontainers.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parse %q: %+v", containerId, err)
			}

			id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing site recovery protection container mapping (%s): %+v", parsedContainerId, err)
					}
				}

				if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
					return tf.ImportAsExistsError("azurerm_site_recovery_replication_policy_vmware_association", *existing.Model.Id)
				}
			}

			parameters := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
					TargetProtectionContainerId: pointer.To(SiteRecoveryReplicationPolicyVMWareAssociationTargetContainerId),
					PolicyId:                    &model.RecoveryReplicationPolicyId,
					ProviderSpecificInput:       replicationprotectioncontainermappings.BaseReplicationProviderSpecificContainerMappingInputImpl{},
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
