// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_profile

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorProfileCreate,
		Read:   resourceCdnFrontDoorProfileRead,
		Update: resourceCdnFrontDoorProfileUpdate,
		Delete: resourceCdnFrontDoorProfileDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := profiles.ParseProfileID(id)
			return err
		}),

		Schema: cdnFrontdoorProfileSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// Verify that they are not downgrading the service from Premium SKU -> Standard SKU...
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				oSku, nSku := diff.GetChange("sku_name")

				if oSku != "" {
					if oSku.(string) == string(profiles.SkuNamePremiumAzureFrontDoor) && nSku.(string) == string(profiles.SkuNameStandardAzureFrontDoor) {
						return fmt.Errorf("downgrading `sku_name` from `%s` to `%s` is not supported", profiles.SkuNamePremiumAzureFrontDoor, profiles.SkuNameStandardAzureFrontDoor)
					}
				}

				return nil
			}),
		),
	}
}
