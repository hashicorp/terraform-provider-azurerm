// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := channels.ParseChannelID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.EventGrid.Channels.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *channels.ChannelId, model *channels.Channel) error {
	state := EventGridPartnerNamespaceChannelTopicResourceModel{
		Name:               id.ChannelName,
		PartnerNamespaceId: channels.NewPartnerNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.PartnerNamespaceName).ID(),
	}
	if model != nil {
		if props := model.Properties; props != nil {
			if props.ChannelType == nil || *props.ChannelType != channels.ChannelTypePartnerTopic {
				return fmt.Errorf("%s is not a Partner Topic channel", *id)
			}
			state.ActivationMessage = pointer.From(props.MessageForActivation)
			state.ExpirationTimeIfNotActivated = pointer.From(props.ExpirationTimeIfNotActivatedUtc)
			state.ReadinessState = pointer.FromEnum(props.ReadinessState)
			if topic := props.PartnerTopicInfo; topic != nil {
				state.PartnerTopic = []ChannelPartnerTopicModel{{
					Name:              pointer.From(topic.Name),
					ResourceGroupName: pointer.From(topic.ResourceGroupName),
					SubscriptionId:    pointer.From(topic.AzureSubscriptionId),
					Source:            pointer.From(topic.Source),
				}}
				state.EventTypes = flattenChannelTopicEventTypes(topic.EventTypeInfo)
			}
		}
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func flattenChannelTopicEventTypes(input *channels.EventTypeInfo) []ChannelTopicEventTypeModel {
	output := make([]ChannelTopicEventTypeModel, 0)
	if input == nil || input.InlineEventTypes == nil {
		return output
	}
	for name, event := range *input.InlineEventTypes {
		output = append(output, ChannelTopicEventTypeModel{
			Name:             name,
			DataSchemaUrl:    pointer.From(event.DataSchemaURL),
			Description:      pointer.From(event.Description),
			DisplayName:      pointer.From(event.DisplayName),
			DocumentationUrl: pointer.From(event.DocumentationURL),
		})
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Name < output[j].Name })
	return output
}
