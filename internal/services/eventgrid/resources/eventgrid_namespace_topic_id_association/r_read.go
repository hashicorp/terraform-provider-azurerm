// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient_v2025_02_15
			var config EventGridNamespaceTopicIdAssociationModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EventGridNamespaceTopicIdAssociationModel{}
			topicSpacesConfigurationExisted := false

			if model := resp.Model; model != nil {
				state.EventGridNamespaceId = pointer.From(model.Id)

				if props := model.Properties; props != nil {
					if topicSpacesConfiguration := props.TopicSpacesConfiguration; topicSpacesConfiguration != nil {
						if routeTopicResourceId := pointer.From(topicSpacesConfiguration.RouteTopicResourceId); routeTopicResourceId != "" {
							state.EventGridNamespaceTopicId = routeTopicResourceId
							topicSpacesConfigurationExisted = true
						}
					}
				}
			}

			if !topicSpacesConfigurationExisted {
				return metadata.MarkAsGone(id)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}
