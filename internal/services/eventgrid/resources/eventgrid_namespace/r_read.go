// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := EventGridNamespaceResourceModel{
				Name:          id.NamespaceName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if model.Sku != nil {
					state.Sku = pointer.FromEnum(model.Sku.Name)
					state.Capacity = pointer.From(model.Sku.Capacity)
				}
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(flattenedIdentity)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					topicSpacesConfig, err := flattenTopicSpacesConfiguration(props.TopicSpacesConfiguration)
					if err != nil {
						return fmt.Errorf("flattening `topic_spaces_configuration`: %v", err)
					}
					state.TopicSpacesConfiguration = topicSpacesConfig
					state.InboundIpRules = flattenInboundIPRules(props.InboundIPRules)
					state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
