// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_static_cidr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.StaticCidrs
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerIpamPoolStaticCidrResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ipamPoolId, err := staticcidrs.ParseIPamPoolID(config.IpamPoolId)
			if err != nil {
				return err
			}

			id := staticcidrs.NewStaticCidrID(subscriptionId, ipamPoolId.ResourceGroupName, ipamPoolId.NetworkManagerName, ipamPoolId.IpamPoolName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := staticcidrs.StaticCidr{
				Name:       pointer.To(config.Name),
				Properties: &staticcidrs.StaticCidrProperties{},
			}

			if len(config.AddressPrefixes) > 0 {
				payload.Properties.AddressPrefixes = pointer.To(config.AddressPrefixes)
			}

			if config.NumberOfIPAddressesToAllocate != "" {
				payload.Properties.NumberOfIPAddressesToAllocate = pointer.To(config.NumberOfIPAddressesToAllocate)
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
