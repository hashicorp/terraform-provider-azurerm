// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient
			namespacesClient := metadata.Client.EventGrid.NamespacesClient

			var model EventGridNamespaceTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			namespaceId, err := namespaces.ParseNamespaceID(model.EventgridNamespaceId)
			if err != nil {
				return err
			}

			id := namespacetopics.NewNamespaceTopicID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			namespaceTopic := namespacetopics.NamespaceTopic{
				Name: pointer.To(model.Name),
				Properties: &namespacetopics.NamespaceTopicProperties{
					EventRetentionInDays: pointer.To(model.EventRetentionInDays),
					InputSchema:          pointer.To(namespacetopics.EventInputSchemaCloudEventSchemaVOneZero),
					PublisherType:        pointer.To(namespacetopics.PublisherTypeCustom),
				},
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, namespaceTopic, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if model.RoutingEnabled {
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

				if existing.Model.Properties.TopicSpacesConfiguration == nil {
					existing.Model.Properties.TopicSpacesConfiguration = &namespaces.TopicSpacesConfiguration{}
				}

				existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = pointer.To(id.ID())

				if err := namespacesClient.CreateOrUpdateThenPoll(ctx, *namespaceId, *existing.Model); err != nil {
					return fmt.Errorf("updating Namespace Topic ID for %s: %+v", *namespaceId, err)
				}
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
