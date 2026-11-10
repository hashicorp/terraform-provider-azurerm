// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels
			id, err := channels.ParseChannelID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var config EventGridPartnerNamespaceChannelTopicResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", *id)
			}
			props := existing.Model.Properties
			if props.ChannelType == nil || *props.ChannelType != channels.ChannelTypePartnerTopic || props.PartnerTopicInfo == nil {
				return fmt.Errorf("%s is not a Partner Topic channel", *id)
			}
			if metadata.ResourceData.HasChange("event_type") {
				props.PartnerTopicInfo.EventTypeInfo = expandChannelTopicEventTypes(config.EventTypes)
			}

			// The default activation exceeds the API's max length validation of 256 and causes problem if included
			// https://github.com/Azure/azure-rest-api-specs/issues/37280
			props.MessageForActivation = nil

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
