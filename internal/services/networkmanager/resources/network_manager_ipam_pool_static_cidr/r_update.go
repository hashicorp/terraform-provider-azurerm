// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_static_cidr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.StaticCidrs

			id, err := staticcidrs.ParseStaticCidrID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerIpamPoolStaticCidrResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", id)
			}

			parameters := resp.Model

			if metadata.ResourceData.HasChange("address_prefixes") {
				if len(model.AddressPrefixes) > 0 {
					parameters.Properties.AddressPrefixes = pointer.To(model.AddressPrefixes)
					// Set nil for AddressPrefixes when changing from `NumberOfIPAddressesToAllocate` to `AddressPrefixes` but the change for `NumberOfIPAddressesToAllocate` is not detected due to the diffSuppressFunc.
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To("")
				} else {
					parameters.Properties.AddressPrefixes = pointer.To([]string{})
				}
			}

			if metadata.ResourceData.HasChange("number_of_ip_addresses_to_allocate") {
				if model.NumberOfIPAddressesToAllocate != "" {
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To(model.NumberOfIPAddressesToAllocate)
					parameters.Properties.AddressPrefixes = pointer.To([]string{})
				} else {
					parameters.Properties.NumberOfIPAddressesToAllocate = pointer.To("")
				}
			}

			if _, err := client.Create(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
