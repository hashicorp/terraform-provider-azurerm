// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient_v2025_02_15
			var config EventGridNamespaceTopicIdAssociationModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namespaces.ParseNamespaceID(config.EventGridNamespaceId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if existing.Model.Properties.TopicSpacesConfiguration == nil {
				existing.Model.Properties.TopicSpacesConfiguration = &namespaces.TopicSpacesConfiguration{}
			} else if rawRouteTopicResourceId := existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId; rawRouteTopicResourceId != nil && !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				routeTopicResourceId := pointer.From(rawRouteTopicResourceId)
				if _, err := namespacetopics.ParseNamespaceTopicID(routeTopicResourceId); err == nil {
					// If there is existing namespace topic ID assigned to `azurerm_eventgrid_namespace`, tell user to import it
					return tf.ImportAsExistsAssociationError("azurerm_eventgrid_namespace_topic_id_association", config.EventGridNamespaceId, routeTopicResourceId)
				}
			}

			existing.Model.Properties.TopicSpacesConfiguration.RouteTopicResourceId = pointer.To(config.EventGridNamespaceTopicId)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, *id, *existing.Model, metadata.SetIDAndIdentityCallback(id)); err != nil {
				return fmt.Errorf("updating Namespace Topic ID for %s: %+v", *id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}
