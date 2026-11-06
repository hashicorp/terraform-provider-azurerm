// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient_v2025_02_15
			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// If parent resource is not found or `TopicSpacesConfiguration` property cannot be accessed, assume that `RouteTopicResourceId` is already unset and continue to remove this resource from state
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return nil
			}
			if resp.Model.Properties == nil {
				return nil
			}

			if resp.Model.Properties.TopicSpacesConfiguration == nil {
				return nil
			}

			resp.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = nil
			if err := client.CreateOrUpdateThenPoll(ctx, *id, *resp.Model); err != nil {
				return fmt.Errorf("removing Namespace Topic ID from %s: %+v", *id, err)
			}

			return nil
		},
	}
}
