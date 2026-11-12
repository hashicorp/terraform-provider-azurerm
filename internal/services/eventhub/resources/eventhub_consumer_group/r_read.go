// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_consumer_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Eventhub.ConsumerGroupClient
			id, err := consumergroups.ParseConsumerGroupID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
		Timeout: 5 * time.Minute,
	}
}

type ConsumerGroupObject struct {
	Name              string `tfschema:"name"`
	NamespaceName     string `tfschema:"namespace_name"`
	EventHubName      string `tfschema:"eventhub_name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	UserMetadata      string `tfschema:"user_metadata"`
}
