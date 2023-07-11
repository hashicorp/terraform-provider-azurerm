// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const SiteRecoveryReplicationPolicyVMWareAssociationTargetContainerId string = "Microsoft Azure"

type SiteRecoveryReplicationPolicyVmwareAssociationModel struct {
	Name                        string `tfschema:"name"`
	RecoveryVaultId             string `tfschema:"recovery_vault_id"`
	RecoveryReplicationPolicyId string `tfschema:"policy_id"`
}

type VMWareReplicationPolicyAssociationResource struct{}

var _ sdk.Resource = VMWareReplicationPolicyAssociationResource{}

func (s VMWareReplicationPolicyAssociationResource) ModelObject() interface{} {
	return &SiteRecoveryReplicationPolicyVmwareAssociationModel{}
}

func (s VMWareReplicationPolicyAssociationResource) ResourceType() string {
	return "azurerm_site_recovery_vmware_replication_policy_association"
}

func (s VMWareReplicationPolicyAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},

		"policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},
	}
}

func (s VMWareReplicationPolicyAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s VMWareReplicationPolicyAssociationResource) Create() sdk.ResourceFunc {
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

			err = validateReplicationPolicyAssociationVaultConfig(ctx, settingsClient, model.RecoveryVaultId)
			if err != nil {
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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing site recovery protection container mapping (%s): %+v", parsedContainerId, err)
				}
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return tf.ImportAsExistsError("azurerm_site_recovery_replication_policy_vmware_association", *existing.Model.Id)
			}

			type RawProviderSpecificInput struct {
				Type   string                 `json:"-"`
				Values map[string]interface{} `json:"-"`
			}

			parameters := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
					TargetProtectionContainerId: utils.String(SiteRecoveryReplicationPolicyVMWareAssociationTargetContainerId),
					PolicyId:                    &model.RecoveryReplicationPolicyId,
					ProviderSpecificInput:       &RawProviderSpecificInput{},
				},
			}

			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (s VMWareReplicationPolicyAssociationResource) Read() sdk.ResourceFunc {
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

func (s VMWareReplicationPolicyAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient

			input := replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.RemoveProtectionContainerMappingInputProperties{
					ProviderSpecificInput: &replicationprotectioncontainermappings.ReplicationProviderContainerUnmappingInput{},
				},
			}

			err = client.DeleteThenPoll(ctx, *id, input)
			if err != nil {
				return fmt.Errorf("deleting %s : %+v", id.String(), err)
			}

			return nil
		},
	}
}

func (s VMWareReplicationPolicyAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ReplicationProtectionContainerMappingsID
}

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

	return handleAzureSdkForGoBug2824(*resp.Items[0].Id), nil
}
