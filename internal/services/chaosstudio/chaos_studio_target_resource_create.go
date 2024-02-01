// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chaosstudio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targettypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r ChaosStudioTargetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Targets
			targetTypesClient := metadata.Client.ChaosStudio.V20231101.TargetTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ChaosStudioTargetResourceSchema

			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			locationId := targettypes.NewLocationID(subscriptionId, config.Location)
			resp, err := targetTypesClient.ListComplete(ctx, locationId, targettypes.DefaultListOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving list of chaos target types: %+v", err)
			}

			// Validate name which needs to be of format [publisher]-[targetType]
			// Only certain target types are accepted, and they could vary by region
			targetTypes := map[string][]string{}
			nameIsValid := false
			for _, item := range resp.Items {
				if name := item.Name; name != nil {
					if strings.EqualFold(config.TargetType, *item.Name) {
						nameIsValid = true
					}
					targetTypes[*item.Name] = pointer.From(item.Properties.ResourceTypes)
				}
			}

			if !nameIsValid {
				return fmt.Errorf("%q is not a valid `target_type` for the region %s, must be one of %+v", config.TargetType, config.Location, targetTypes)
			}

			id := commonids.NewChaosStudioTargetID(config.TargetResourceId, config.TargetType)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload targets.Target

			// The API only accepts requests with an empty body for Properties
			props := struct{}{}

			payload.Location = pointer.To(location.Normalize(config.Location))
			payload.Properties = pointer.To(props)

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
