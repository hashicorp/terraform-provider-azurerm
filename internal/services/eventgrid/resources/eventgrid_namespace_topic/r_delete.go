// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient
			namespacesClient := metadata.Client.EventGrid.NamespacesClient

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}

			namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
			existing, err := namespacesClient.Get(ctx, namespaceId)
			if err == nil && existing.Model != nil && existing.Model.Properties != nil {
				if topicSpacesConfiguration := existing.Model.Properties.TopicSpacesConfiguration; topicSpacesConfiguration != nil {
					if routeTopicResourceId := topicSpacesConfiguration.RouteTopicResourceId; routeTopicResourceId != nil && strings.EqualFold(id.ID(), pointer.From(routeTopicResourceId)) {
						existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = nil
						if err := namespacesClient.CreateOrUpdateThenPoll(ctx, namespaceId, *existing.Model); err != nil {
							return fmt.Errorf("updating Namespace Topic ID for %s: %+v", namespaceId, err)
						}
					}
				}
			}

			return nil
		},
	}
}
