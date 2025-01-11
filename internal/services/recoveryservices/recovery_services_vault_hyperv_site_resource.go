// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HyperVSiteModel struct {
	Name            string `tfschema:"name"`
	RecoveryVaultId string `tfschema:"recovery_vault_id"`
}

type HyperVSiteResource struct{}

var _ sdk.Resource = HyperVSiteResource{}

func (r HyperVSiteResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			ValidateFunc: replicationfabrics.ValidateVaultID,
		},
	}
}

func (r HyperVSiteResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r HyperVSiteResource) ModelObject() interface{} {
	return &HyperVSiteModel{}
}

func (r HyperVSiteResource) ResourceType() string {
	return "azurerm_site_recovery_services_vault_hyperv_site"
}

func (r HyperVSiteResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel HyperVSiteModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %s", err)
			}

			client := metadata.Client.RecoveryServices.FabricClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, err := replicationfabrics.ParseVaultID(metaModel.RecoveryVaultId)
			if err != nil {
				return err
			}

			id := replicationfabrics.NewReplicationFabricID(subscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, metaModel.Name)

			// the instance type `HyperVSite` is not exposed in Swagger, tracked on https://github.com/Azure/azure-rest-api-specs/issues/22016
			parameters := replicationfabrics.FabricCreationInput{
				Properties: &replicationfabrics.FabricCreationInputProperties{
					CustomDetails: replicationfabrics.BaseFabricSpecificCreationInputImpl{
						InstanceType: "HyperVSite",
					},
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

func (r HyperVSiteResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.FabricClient
			id, err := replicationfabrics.ParseReplicationFabricID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id, replicationfabrics.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			state := HyperVSiteModel{
				Name: id.ReplicationFabricName,
			}

			vaultId := replicationfabrics.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
			state.RecoveryVaultId = vaultId.ID()

			return metadata.Encode(&state)
		},
	}
}

func (r HyperVSiteResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute, // when a host connected to site, it will cost up to 180 minutes to delete.
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.FabricClient
			id, err := replicationfabrics.ParseReplicationFabricID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r HyperVSiteResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationfabrics.ValidateReplicationFabricID
}
