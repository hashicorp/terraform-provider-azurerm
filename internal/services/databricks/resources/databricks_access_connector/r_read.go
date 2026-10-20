// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_access_connector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.DataBricks.AccessConnectorClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
			}

			state := AccessConnectorResourceModel{
				Name:          id.AccessConnectorName,
				Location:      location.NormalizeNilable(pointer.To(resp.Model.Location)),
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if model.Identity != nil {
					identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}

					if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
						return fmt.Errorf("setting `identity`: %+v", err)
					}
				}
			}
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}
