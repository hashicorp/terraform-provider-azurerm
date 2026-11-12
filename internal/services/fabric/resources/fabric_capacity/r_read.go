// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Fabric.FabricCapacitiesClient

			id, err := fabriccapacities.ParseCapacityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := FabricCapacityResourceModel{
				Name:              id.CapacityName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.AdministrationMembers = model.Properties.Administration.Members
				state.Sku = flattenSkuModel(model.Sku)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenSkuModel(input fabriccapacities.RpSku) []SkuModel {
	outputList := make([]SkuModel, 0)
	output := SkuModel{
		Name: input.Name,
		Tier: string(input.Tier),
	}

	return append(outputList, output)
}
