// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipampools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.IPamPools
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerIpamPoolResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := ipampools.ParseNetworkManagerID(config.NetworkManagerId)
			if err != nil {
				return err
			}

			id := ipampools.NewIPamPoolID(subscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := ipampools.IPamPool{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				Properties: ipampools.IPamPoolProperties{
					AddressPrefixes: config.AddressPrefixes,
					Description:     pointer.To(config.Description),
					DisplayName:     pointer.To(config.DisplayName),
					ParentPoolName:  pointer.To(config.ParentPoolName),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, payload, ipampools.DefaultCreateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
