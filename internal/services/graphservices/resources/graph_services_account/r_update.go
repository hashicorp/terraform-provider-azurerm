// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			id, err := graphservicesprods.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AccountResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.AccountsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: model was nil", *id)
			}

			payload := graphservicesprods.TagUpdate{}
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if _, err := client.AccountsUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
