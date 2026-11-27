// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_accelerator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).ApplicationAcceleratorClient

			// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
			// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
			id, err := appplatform.ParseApplicationAcceleratorIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudAcceleratorModel{
				Name:                 id.ApplicationAcceleratorName,
				SpringCloudServiceId: commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName).ID(),
			}
			return metadata.Encode(&state)
		},
	}
}
