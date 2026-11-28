// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			searchServiceId := sharedprivatelinkresources.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName)
			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			var state SharedPrivateLinkServiceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient

			if metadata.ResourceData.HasChange("request_message") {
				props := sharedprivatelinkresources.SharedPrivateLinkResource{
					Properties: &sharedprivatelinkresources.SharedPrivateLinkResourceProperties{
						RequestMessage: pointer.To(state.RequestMessage),
					},
				}
				if err := client.CreateOrUpdateThenPoll(ctx, *id, props, sharedprivatelinkresources.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
		Timeout: 60 * time.Minute,
	}
}
