// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const TargetContainerIdAzure = "Microsoft Azure"

// we only support replicate to Azure as backup to customer-managed sites is being deprecated.
// https://learn.microsoft.com/en-us/azure/site-recovery/site-to-site-deprecation

type HyperVReplicationPolicyAssociationModel struct {
	Name     string `tfschema:"name"`
	FabricId string `tfschema:"hyperv_site_id"`
	PolicyId string `tfschema:"policy_id"`
}

type HyperVReplicationPolicyAssociationResource struct{}

var _ sdk.Resource = HyperVReplicationPolicyAssociationResource{}

func (h HyperVReplicationPolicyAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hyperv_site_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},
	}
}

func (h HyperVReplicationPolicyAssociationResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (h HyperVReplicationPolicyAssociationResource) ModelObject() interface{} {
	return &HyperVReplicationPolicyAssociationModel{}
}

func (h HyperVReplicationPolicyAssociationResource) ResourceType() string {
	return "azurerm_site_recovery_hyperv_replication_policy_association"
}

func (h HyperVReplicationPolicyAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationprotectioncontainermappings.ValidateReplicationProtectionContainerMappingID
}

func (h HyperVReplicationPolicyAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan HyperVReplicationPolicyAssociationModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			parsedFabricId, err := replicationfabrics.ParseReplicationFabricID(plan.FabricId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.FabricId, err)
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			containerId, err := fetchHyperVReplicationPolicyAssociationContainerNameByHostName(ctx, metadata.Client.RecoveryServices.ProtectionContainerClient, *parsedFabricId)
			if err != nil {
				return fmt.Errorf("fetching container id: %+v", err)
			}

			parsedContainerId, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", containerId, err)
			}

			id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID(subscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, plan.Name)

			type hyperVMappingSpecificInput struct { // a workaround for https://github.com/Azure/azure-rest-api-specs/issues/22769
				InstanceType string `json:"instanceType"`
			}

			param := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
					PolicyId:                    &plan.PolicyId,
					TargetProtectionContainerId: utils.String(TargetContainerIdAzure),
					ProviderSpecificInput:       hyperVMappingSpecificInput{},
				},
			}
			if err := client.CreateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (h HyperVReplicationPolicyAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ContainerMappingClient

			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			fabricId := replicationfabrics.NewReplicationFabricID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := HyperVReplicationPolicyAssociationModel{
				Name:     id.ReplicationProtectionContainerMappingName,
				FabricId: fabricId.ID(),
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `Properties` was nil", id)
			}

			prop := resp.Model.Properties
			if prop.PolicyId != nil {
				state.PolicyId = handleAzureSdkForGoBug2824(*prop.PolicyId)
			}

			return metadata.Encode(&state)
		},
	}
}

func (h HyperVReplicationPolicyAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ContainerMappingClient

			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			if err := client.DeleteThenPoll(ctx, *id, replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func fetchHyperVReplicationPolicyAssociationContainerNameByHostName(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, fabricId replicationfabrics.ReplicationFabricId) (string, error) {
	id, err := replicationprotectioncontainers.ParseReplicationFabricID(fabricId.ID())
	if err != nil {
		return "", fmt.Errorf("parsing %s: %+v", fabricId.ID(), err)
	}

	resp, err := containerClient.ListByReplicationFabricsComplete(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("listing containers: %+v", err)
	}

	if len(resp.Items) == 0 || len(resp.Items) > 1 {
		return "", fmt.Errorf("expected one container but got %d", len(resp.Items))
	}

	if resp.Items[0].Id == nil {
		return "", fmt.Errorf("container id is nil")
	}

	return handleAzureSdkForGoBug2824(*resp.Items[0].Id), nil
}
