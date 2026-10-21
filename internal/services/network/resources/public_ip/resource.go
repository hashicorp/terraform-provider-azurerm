// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpCreate,
		Read:   resourcePublicIpRead,
		Update: resourcePublicIpUpdate,
		Delete: resourcePublicIpDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.PublicIPAddressId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.PublicIPAddressId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: publicIpSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				sku := d.Get("sku").(string)
				if strings.EqualFold(sku, string(publicipaddresses.PublicIPAddressSkuNameBasic)) && d.HasChanges(
					"name",
					"resource_group_name",
					"location",
					"allocation_method",
					"edge_zone",
					"ip_version",
					"sku",
					"sku_tier",
					"public_ip_prefix_id",
					"ip_tags",
					"zones",
				) {
					return errors.New(publicIPBasicSkuCreateDeprecationMessage)
				}

				return nil
			}),
			pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				skuTier := d.Get("sku_tier").(string)
				sku := d.Get("sku").(string)
				if strings.EqualFold(skuTier, string(publicipaddresses.PublicIPAddressSkuTierGlobal)) && !strings.EqualFold(sku, string(publicipaddresses.PublicIPAddressSkuNameStandard)) {
					return errors.New("`sku` must be set to `Standard` when `sku_tier` is set to `Global`")
				}
				return nil
			}),
			pluginsdk.ForceNewIfChange("domain_name_label_scope", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != "" || new.(string) == ""
			}),
		),
	}
}
