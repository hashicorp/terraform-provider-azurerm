// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.Tenants

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := AadB2cDirectoryModel{
				DomainName:    id.DirectoryName,
				ResourceGroup: id.ResourceGroup,
				CountryCode:   metadata.ResourceData.Get("country_code").(string),
				DisplayName:   metadata.ResourceData.Get("display_name").(string),
			}

			if model.Location != nil {
				state.DataResidencyLocation = string(*model.Location)
			}

			if model.Sku != nil {
				state.Sku = string(model.Sku.Name)
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if properties := model.Properties; properties != nil {
				if billingConfig := properties.BillingConfig; billingConfig != nil {
					if billingConfig.BillingType != nil {
						state.BillingType = string(*billingConfig.BillingType)
					}
					if billingConfig.EffectiveStartDateUtc != nil {
						state.EffectiveStartDate = *billingConfig.EffectiveStartDateUtc
					}
				}

				if properties.TenantId != nil {
					state.TenantId = *properties.TenantId
				}
			}

			return metadata.Encode(&state)
		},
	}
}
