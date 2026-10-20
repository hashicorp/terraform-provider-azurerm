// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_profile

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_profile/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnProfileCreate,
		Read:   resourceCdnProfileRead,
		Update: resourceCdnProfileUpdate,
		Delete: resourceCdnProfileDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnProfileV0ToV1{},
			// v1 -> v2 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			1: migration.CdnProfileV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := profiles.ParseProfileID(id)
			return err
		}),

		Schema: cdnProfileSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			sku := d.Get("sku").(string)

			if helpers.IsCdnFullyRetired() {
				return fmt.Errorf("%s", helpers.FullyRetiredMessage)
			}

			switch sku {
			case string(cdn.SkuNameStandardAkamai):
				// The Akamai sku was retired on 31 October 2023...
				if helpers.IsCdnStandardAkamaiDeprecatedForCreation() {
					return fmt.Errorf("%s", helpers.AkamaiDeprecationMessage)
				}
			case string(cdn.SkuNameStandardVerizon), string(cdn.SkuNamePremiumVerizon):
				// The Verizon skus were retired on 15 January 2025...
				if helpers.IsCdnStandardVerizonPremiumVerizonDeprecatedForCreation() {
					return fmt.Errorf("%s", helpers.VerizonDeprecationMessage)
				}
			case string(cdn.SkuNameStandardMicrosoft), string(cdn.SkuNameStandardChinaCdn):
				// The only currently supported sku's for CDN are 'Standard_Microsoft' and 'Standard_ChinaCdn' which
				// will also be deprecated as of October 1, 2025, but can still be updated until September 30, 2027.
				// First, check to see if any of the force new fields have changed after the October 1, 2025 deprecation date,
				// if they have we need to block the update because the force new will cause the deployments to fail...
				if helpers.IsCdnDeprecatedForCreation() && d.HasChanges("name", "location", "resource_group_name", "sku") {
					return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
				}
			}

			return nil
		}),
	}
}
