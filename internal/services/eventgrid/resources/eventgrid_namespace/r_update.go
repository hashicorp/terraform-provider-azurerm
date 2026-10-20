// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient

			var model EventGridNamespaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := namespaces.NamespaceUpdateParameters{
				Properties: &namespaces.NamespaceUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = identity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("capacity") {
				payload.Sku = &namespaces.NamespaceSku{
					Capacity: pointer.To(model.Capacity),
				}
			}

			if metadata.ResourceData.HasChange("inbound_ip_rule") {
				payload.Properties.InboundIPRules = expandInboundIPRules(model.InboundIpRules)
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.ToEnum[namespaces.PublicNetworkAccess](model.PublicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("topic_spaces_configuration") {
				payload.Properties.TopicSpacesConfiguration = expandTopicSpacesConfigurationUpdate(model.TopicSpacesConfiguration)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", *id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
