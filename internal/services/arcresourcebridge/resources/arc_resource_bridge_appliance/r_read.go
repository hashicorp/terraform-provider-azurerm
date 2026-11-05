// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.ArcResourceBridge.AppliancesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %q: %+v", *id, err)
			}

			state := ApplianceModel{
				Name:              id.ApplianceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Identity = identity.FlattenSystemAssignedToModel(model.Identity)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.Distro = pointer.From(props.Distro)
					state.PublicKeyBase64 = pointer.From(props.PublicKey)

					if infraConfig := props.InfrastructureConfig; infraConfig != nil {
						state.Provider = pointer.From(infraConfig.Provider)
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}
