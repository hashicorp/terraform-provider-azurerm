// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bastion_host

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/bastionhosts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBastionHostCreate,
		Read:   resourceBastionHostRead,
		Update: resourceBastionHostUpdate,
		Delete: resourceBastionHostDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&bastionhosts.BastionHostId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&bastionhosts.BastionHostId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: bastionHostSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku", func(ctx context.Context, old, new, meta any) bool {
				// downgrade the SKU is not supported, recreate the resource
				if old.(string) != "" && new.(string) != "" {
					return skuWeight[old.(string)] > skuWeight[new.(string)]
				}
				return false
			}),
			func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
				sku := bastionhosts.BastionHostSkuName(d.Get("sku").(string))

				// GetRawConfig is used because `public_ip_address_id` may reference another resource and be unknown during plan.
				// d.Get() returns "" for unknown values, which would incorrectly trigger the required check.
				// By inspecting the raw config, we can distinguish between "not set" and "unknown" and skip validation when unknown.
				rawConfig := d.GetRawConfig()
				ipConfigRaw := rawConfig.AsValueMap()["ip_configuration"]

				// Basic and Standard SKUs require public_ip_address_id; Developer uses virtual_network_id instead.
				// Premium without public_ip_address_id enables private-only mode.
				if sku != bastionhosts.BastionHostSkuNamePremium && sku != bastionhosts.BastionHostSkuNameDeveloper {
					if ipConfigRaw.IsNull() || !ipConfigRaw.IsKnown() || len(ipConfigRaw.AsValueSlice()) == 0 {
						return errors.New("`ip_configuration` with `public_ip_address_id` is required when `sku` is `Basic` or `Standard`")
					}
					pipRaw := ipConfigRaw.AsValueSlice()[0].AsValueMap()["public_ip_address_id"]
					if pipRaw.IsKnown() && (pipRaw.IsNull() || pipRaw.AsString() == "") {
						return errors.New("`public_ip_address_id` must be set in `ip_configuration` when `sku` is `Basic` or `Standard`")
					}
				}

				return nil
			},
		),
	}
}
