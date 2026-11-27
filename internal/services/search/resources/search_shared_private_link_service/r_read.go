// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := searchClient.From(metadata.Client).SearchSharedPrivateLinkResourceClient
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, sharedprivatelinkresources.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &SharedPrivateLinkServiceModel{
				Name:            id.SharedPrivateLinkResourceName,
				SearchServiceId: services.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.GroupId != nil {
						state.SubResourceName = *props.GroupId
					}

					if props.PrivateLinkResourceId != nil {
						state.TargetResourceId = *props.PrivateLinkResourceId
					}

					if props.RequestMessage != nil {
						state.RequestMessage = *props.RequestMessage
					}

					if props.Status != nil {
						state.Status = string(*props.Status)
					}
				}
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}
