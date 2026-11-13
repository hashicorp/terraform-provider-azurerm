// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.Tenants

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AadB2cDirectoryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			properties := tenants.UpdateTenant{
				Sku: tenants.Sku{
					Name: tenants.SkuName(state.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &state.Tags,
			}

			if _, err := client.Update(ctx, *id, properties); err != nil {
				return err
			}

			return nil
		},
	}
}
