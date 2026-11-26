// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient

			var model EventGridNamespaceTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := namespacetopics.NamespaceTopicUpdateParameters{
				Properties: &namespacetopics.NamespaceTopicUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("event_retention_in_days") {
				payload.Properties.EventRetentionInDays = pointer.To(model.EventRetentionInDays)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
