// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_management_private_link

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceManagementPrivateLinkClient

			id, err := resourcemanagementprivatelink.ParseResourceManagementPrivateLinkID(metadata.ResourceData.Id())
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

			schema := ResourceManagementPrivateLinkResourceSchema{
				Name:              id.ResourceManagementPrivateLinkName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.NormalizeNilable(model.Location)
			}

			return metadata.Encode(&schema)
		},
	}
}

type ResourceManagementPrivateLinkResourceSchema struct {
	Location          string `tfschema:"location"`
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
}
