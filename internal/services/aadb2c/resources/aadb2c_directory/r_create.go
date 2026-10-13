// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.Tenants
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AadB2cDirectoryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.CountryCode == "" {
				return errors.New("`country_code` is required when creating a new AADB2C directory")
			}
			if model.DisplayName == "" {
				return errors.New("`display_name` is required when creating a new AADB2C directory")
			}

			id := tenants.NewB2CDirectoryID(subscriptionId, model.ResourceGroup, model.DomainName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			availabilityResult, err := client.CheckNameAvailability(ctx, commonids.NewSubscriptionID(subscriptionId), tenants.CheckNameAvailabilityRequest{
				Name:        &model.DomainName,
				CountryCode: &model.CountryCode,
			})
			if err != nil {
				return fmt.Errorf("checking availability of `domain_name`: %v", err)
			}

			if availabilityResult.Model != nil {
				if availabilityResult.Model.NameAvailable == nil || !*availabilityResult.Model.NameAvailable {
					reason := "unknown reason"
					if availabilityResult.Model.Reason != nil {
						reason = *availabilityResult.Model.Reason
					}
					if availabilityResult.Model.Message != nil {
						reason = fmt.Sprintf("%s (%s)", reason, *availabilityResult.Model.Message)
					}
					return fmt.Errorf("checking availability of `domain_name`: the specified domain %q is unavailable: %s", model.DomainName, reason)
				}
			}

			properties := tenants.CreateTenant{
				Location: tenants.Location(model.DataResidencyLocation),
				Properties: tenants.TenantPropertiesForCreate{
					CreateTenantProperties: tenants.CreateTenantProperties{
						CountryCode: model.CountryCode,
						DisplayName: model.DisplayName,
					},
				},
				Sku: tenants.Sku{
					Name: tenants.SkuName(model.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateCallbackThenPoll(ctx, id, properties, metadata.SetIDCallback(&id)); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}
