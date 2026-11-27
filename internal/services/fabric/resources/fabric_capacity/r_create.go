// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	fabricClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/fabric/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := fabricClient.From(metadata.Client).FabricCapacitiesClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var model FabricCapacityResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fabriccapacities.NewCapacityID(subscriptionId, model.ResourceGroupName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if len(model.AdministrationMembers) == 0 {
				return errors.New("`administration_members` is required when creating a fabric capacity")
			}

			properties := fabriccapacities.FabricCapacity{
				Location: location.Normalize(model.Location),
				Properties: fabriccapacities.FabricCapacityProperties{
					Administration: fabriccapacities.CapacityAdministration{
						Members: model.AdministrationMembers,
					},
				},
				Sku: expandSkuModel(model.Sku),
			}

			if model.Tags != nil {
				properties.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandSkuModel(inputList []SkuModel) fabriccapacities.RpSku {
	input := &inputList[0]
	return fabriccapacities.RpSku{
		Name: input.Name,
		Tier: fabriccapacities.RpSkuTier(input.Tier),
	}
}
