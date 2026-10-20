// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_consumer_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := consumergroups.ParseConsumerGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ConsumerGroupObject
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Eventhub.ConsumerGroupClient

			parameters := consumergroups.ConsumerGroup{
				Name: pointer.To(id.ConsumerGroupName),
				Properties: &consumergroups.ConsumerGroupProperties{
					UserMetadata: pointer.To(state.UserMetadata),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
