// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions

			id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(metadata.ResourceData.Id())
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

			state := DevCenterDevBoxDefinitionResourceModel{
				Name:        id.DevBoxDefinitionName,
				DevCenterId: devboxdefinitions.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if v := props.ImageReference; v != nil {
						state.ImageReferenceId = pointer.From(v.Id)
					}

					state.HibernateSupport = pointer.From(props.HibernateSupport) == devboxdefinitions.HibernateSupportEnabled

					if v := props.Sku; v != nil {
						state.SkuName = flattenDevCenterDevBoxDefinition(props.Sku)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDevCenterDevBoxDefinition(input *devboxdefinitions.Sku) string {
	var skuName string
	if input == nil {
		return skuName
	}

	skuName = pointer.From(input).Name

	return skuName
}
