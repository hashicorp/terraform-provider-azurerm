// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = SiteRecoveryRecoveryVaultDataSource{}

type SiteRecoveryRecoveryVaultDataSource struct{}
type SiteRecoveryRecoveryVaultDataSourceModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Sku               string                                     `tfschema:"sku"`
	Tags              map[string]string                          `tfschema:"tags"`
}

func (SiteRecoveryRecoveryVaultDataSource) ModelObject() interface{} {
	return &SiteRecoveryRecoveryVaultDataSourceModel{}
}

func (SiteRecoveryRecoveryVaultDataSource) ResourceType() string {
	return "azurerm_recovery_services_vault"
}

func (SiteRecoveryRecoveryVaultDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}
func (SiteRecoveryRecoveryVaultDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SiteRecoveryRecoveryVaultDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.VaultsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var recoveryServiceVault SiteRecoveryRecoveryVaultDataSourceModel
			if err := metadata.Decode(&recoveryServiceVault); err != nil {
				return err
			}

			id := vaults.NewVaultID(subscriptionId, recoveryServiceVault.ResourceGroupName, recoveryServiceVault.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)

				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				skuName := ""
				if model.Sku != nil {
					skuName = string(model.Sku.Name)
				}

				recoveryServiceVault.Sku = skuName
				recoveryServiceVault.Location = location.Normalize(model.Location)
				recoveryServiceVault.Tags = pointer.From(model.Tags)
				recoveryServiceVault.Identity = pointer.From(flattenedIdentity)
			}

			metadata.SetID(id)

			if err := metadata.Encode(&recoveryServiceVault); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}
