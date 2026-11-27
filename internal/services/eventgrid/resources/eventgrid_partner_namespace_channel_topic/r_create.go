// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels
			var config EventGridPartnerNamespaceChannelTopicResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parent, err := channels.ParsePartnerNamespaceID(config.PartnerNamespaceId)
			if err != nil {
				return err
			}
			id := channels.NewChannelID(parent.SubscriptionId, parent.ResourceGroupName, parent.PartnerNamespaceName, config.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}
			topic := config.PartnerTopic[0]
			properties := channels.ChannelProperties{
				ChannelType: pointer.To(channels.ChannelTypePartnerTopic),
				PartnerTopicInfo: &channels.PartnerTopicInfo{
					AzureSubscriptionId: pointer.To(topic.SubscriptionId),
					ResourceGroupName:   pointer.To(topic.ResourceGroupName),
					Name:                pointer.To(topic.Name),
					Source:              pointer.To(topic.Source),
					EventTypeInfo:       expandChannelTopicEventTypes(config.EventTypes),
				},
			}
			if _, err := client.CreateOrUpdate(ctx, id, channels.Channel{Properties: &properties}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandChannelTopicEventTypes(input []ChannelTopicEventTypeModel) *channels.EventTypeInfo {
	if len(input) == 0 {
		return nil
	}

	events := make(map[string]channels.InlineEventProperties, len(input))
	for _, event := range input {
		events[event.Name] = channels.InlineEventProperties{
			DataSchemaURL:    pointer.To(event.DataSchemaUrl),
			Description:      pointer.To(event.Description),
			DisplayName:      pointer.To(event.DisplayName),
			DocumentationURL: pointer.To(event.DocumentationUrl),
		}
	}
	return &channels.EventTypeInfo{Kind: pointer.To(channels.EventDefinitionKindInline), InlineEventTypes: &events}
}
