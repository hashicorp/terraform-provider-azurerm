// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := searchClient.From(metadata.Client).SearchSharedPrivateLinkResourceClient
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			searchServiceId := sharedprivatelinkresources.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName)
			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			if err := client.DeleteThenPoll(ctx, *id, sharedprivatelinkresources.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
		Timeout: 60 * time.Minute,
	}
}
