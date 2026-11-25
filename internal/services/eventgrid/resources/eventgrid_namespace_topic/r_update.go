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

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient
			namespacesClient := metadata.Client.EventGrid.NamespacesClient

			var model EventGridNamespaceTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := namespacetopics.NamespaceTopicUpdateParameters{
				Properties: &namespacetopics.NamespaceTopicUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("event_retention_in_days") {
				payload.Properties.EventRetentionInDays = pointer.To(model.EventRetentionInDays)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("routing_enabled") {
				namespaceId, err := namespaces.ParseNamespaceID(model.EventgridNamespaceId)
				if err != nil {
					return err
				}

				existing, err := namespacesClient.Get(ctx, *namespaceId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *namespaceId, err)
				}

				if existing.Model == nil {
					return fmt.Errorf("retrieving %s: `model` was nil", *namespaceId)
				}
				if existing.Model.Properties == nil {
					return fmt.Errorf("retrieving %s: `properties` was nil", *namespaceId)
				}

				topicSpacesConfiguration := existing.Model.Properties.TopicSpacesConfiguration
				routeTopicResourceIdUpdated := false

				if model.RoutingEnabled {
					if topicSpacesConfiguration == nil {
						existing.Model.Properties.TopicSpacesConfiguration = &namespaces.TopicSpacesConfiguration{}
					}

					existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = pointer.To(id.ID())
					routeTopicResourceIdUpdated = true
				} else if topicSpacesConfiguration != nil && strings.EqualFold(id.ID(), pointer.From(topicSpacesConfiguration.RouteTopicResourceId)) {
					existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = nil
					routeTopicResourceIdUpdated = true
				}

				if routeTopicResourceIdUpdated {
					if err := namespacesClient.CreateOrUpdateThenPoll(ctx, *namespaceId, *existing.Model); err != nil {
						return fmt.Errorf("updating Namespace Topic ID for %s: %+v", *namespaceId, err)
					}
				}
			}

			return nil
		},
	}
}
